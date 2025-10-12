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

const (
	nackDelay          = 30 * time.Second
	batchMaxWait       = 5 * time.Second
	maxConcurrentPulls = 8
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
			jsCfg.StreamName,
			jsCfg.StreamName + ".>",
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

	subject := streamName + ".>"
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
			_ = msg.NakWithDelay(nackDelay)

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

func (n *NATSMq) SubscribeBatch(
	ctx context.Context,
	streamName string,
	batchSize int,
	processFn func(ctx context.Context, messages [][]byte) error,
) error {
	js, ok := n.streams[streamName]
	if !ok {
		return fmt.Errorf("stream not found: %s", streamName)
	}

	subject := streamName + ".>"
	durable := streamName + "_batch"

	sub, err := js.PullSubscribe(subject, durable,
		nats.BindStream(streamName),
		nats.PullMaxWaiting(maxConcurrentPulls), // max concurrent pulls
	)
	if err != nil {
		return fmt.Errorf("create pull subscriber failed: %w", err)
	}

	slog.Info("nats: pull subscriber started",
		"stream", streamName,
		"subject", subject,
		"batch_size", batchSize,
	)

	for {
		select {
		case <-ctx.Done():
			slog.Info("nats: stopping batch subscriber", "stream", streamName)

			return sub.Drain()

		default:
			// Try to fetch a batch of messages
			msgs, err := sub.Fetch(batchSize, nats.MaxWait(batchMaxWait))
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue // just no messages — not an error
				}

				slog.Error("nats: fetch failed", "err", err)
				time.Sleep(time.Second)

				continue
			}

			// Collect payloads
			payloads := make([][]byte, 0, len(msgs))
			for _, msg := range msgs {
				payloads = append(payloads, msg.Data)
			}

			// Process batch
			if err := processFn(ctx, payloads); err != nil {
				slog.Error("nats: batch processing failed", "stream", streamName, "err", err)
				for _, msg := range msgs {
					_ = msg.NakWithDelay(nackDelay)
				}

				continue
			}

			// Ack all messages
			for _, msg := range msgs {
				if err := msg.Ack(); err != nil {
					slog.Warn("nats: ack failed", "err", err)
				}
			}
		}
	}
}

func (n *NATSMq) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
