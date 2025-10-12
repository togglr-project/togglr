package natsmq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/togglr-project/togglr/internal/domain"
)

type Config struct {
	URL        string
	JetStreams []JetStreamConfig
}

type JetStreamConfig struct {
	StreamName string
	MaxAge     time.Duration
}

type NATSMq struct {
	conn    *nats.Conn
	streams map[string]nats.JetStreamContext
}

// New initializes connection and JetStream streams based on config.
func New(cfg *Config) (*NATSMq, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	mq := &NATSMq{
		conn:    conn,
		streams: make(map[string]nats.JetStreamContext),
	}

	for _, jsCfg := range cfg.JetStreams {
		js, err := conn.JetStream()
		if err != nil {
			return nil, fmt.Errorf("init JetStream: %w", err)
		}

		subjects := []string{
			fmt.Sprintf("%s", jsCfg.StreamName),
			fmt.Sprintf("%s.>", jsCfg.StreamName),
		}

		if _, err := js.StreamInfo(jsCfg.StreamName); err != nil {
			if errors.Is(err, nats.ErrStreamNotFound) {
				slog.Info("nats: creating stream",
					"stream", jsCfg.StreamName,
					"subjects", subjects,
					"max_age", jsCfg.MaxAge,
				)
				_, err = js.AddStream(&nats.StreamConfig{
					Name:      jsCfg.StreamName,
					Subjects:  subjects,
					Retention: nats.LimitsPolicy,
					Storage:   nats.FileStorage,
					MaxAge:    jsCfg.MaxAge,
				})
				if err != nil {
					return nil, fmt.Errorf("create stream %s: %w", jsCfg.StreamName, err)
				}
			} else {
				return nil, fmt.Errorf("get stream info %s: %w", jsCfg.StreamName, err)
			}
		}

		mq.streams[jsCfg.StreamName] = js
	}

	return mq, nil
}

// Publish publishes a message to a JetStream stream.
// subjectSuffix — optional part after the stream name (e.g. ".feature.track").
func (n *NATSMq) Publish(_ context.Context, streamName string, data []byte) error {
	js, ok := n.streams[streamName]
	if !ok {
		return fmt.Errorf("stream not found: %s", streamName)
	}

	subject := streamName + ".event"

	ack, err := js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("publish to stream %q failed: %w", streamName, err)
	}

	slog.Debug("nats: published message",
		"stream", ack.Stream,
		"seq", ack.Sequence,
		"subject", subject,
	)

	return nil
}

// Subscribe creates a queue subscription for a given stream.
func (n *NATSMq) Subscribe(
	ctx context.Context,
	streamName string,
	processFn func(ctx context.Context, data []byte) error,
) error {
	js, ok := n.streams[streamName]
	if !ok {
		return fmt.Errorf("stream not found: %s", streamName)
	}

	subject := fmt.Sprintf("%s.>", streamName)
	queue := streamName

	sub, err := js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("nats: panic in handler", "recover", r)
				_ = msg.Nak()
			}
		}()

		if err := processFn(ctx, msg.Data); err != nil {
			var skippableErr *domain.SkippableError
			if errors.As(err, &skippableErr) {
				slog.Warn("nats: message skipped", "subject", msg.Subject, "err", err)
				_ = msg.Ack()

				return
			}

			slog.Error("nats: message failed", "subject", msg.Subject, "err", err)
			_ = msg.NakWithDelay(30 * time.Second)

			return
		}

		if err := msg.Ack(); err != nil {
			slog.Warn("nats: ack failed", "err", err)
		}
	}, nats.ManualAck(), nats.BindStream(streamName))
	if err != nil {
		return fmt.Errorf("subscribe failed for stream %s: %w", streamName, err)
	}

	defer func() { _ = sub.Drain() }()

	slog.Info("nats: subscribed", "stream", streamName)

	<-ctx.Done()

	return nil
}

func (n *NATSMq) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
