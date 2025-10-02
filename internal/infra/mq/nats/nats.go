package natsmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type NATSMq struct {
	conn *nats.Conn
}

func New(url string) (*NATSMq, error) {
	conn, err := connectNATS(url)
	if err != nil {
		return nil, err
	}

	return &NATSMq{conn: conn}, nil
}

func (n *NATSMq) Publish(ctx context.Context, topic string, data []byte) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := n.conn.Publish(topic, data)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (n *NATSMq) Request(
	ctx context.Context,
	topic string,
	data []byte,
) ([]byte, error) {
	msg, err := n.conn.RequestWithContext(ctx, topic, data)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	resp := msg.Data
	slog.Debug("nats: request response", "topic", topic, "response", string(resp))

	return resp, nil
}

func (n *NATSMq) HandleRequest(
	ctx context.Context,
	topic string,
	handlerFn func(ctx context.Context, data []byte) (response []byte, err error),
) error {
	sub, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		slog.Debug("nats: request handler: process message", "topic", topic, "request", string(msg.Data))

		response, err := handlerFn(ctx, msg.Data)
		if err != nil {
			slog.Error("nats: request handler: process message failed", "error", err)
			errResponse := []byte("error:" + err.Error())
			if err := msg.Respond(errResponse); err != nil {
				slog.Error("nats: failed to send error response", "error", err)
			}

			return
		}

		if err := msg.Respond(response); err != nil {
			slog.Error("nats: failed to send response", "error", err)
		} else {
			slog.Debug("nats: request handler: response sent",
				"topic", topic, "response", string(response))
		}
	})
	if err != nil {
		return fmt.Errorf("subscribe to topic %q failed: %w", topic, err)
	}

	defer func() { _ = sub.Unsubscribe() }()

	<-ctx.Done()

	return nil
}

func (n *NATSMq) Subscribe(
	ctx context.Context,
	topic string,
	processFn func(ctx context.Context, data []byte) error,
) error {
	sub, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		err := processFn(ctx, msg.Data)
		if err != nil {
			slog.Error("nats: failed to process message", "error", err)
		}
	})
	if err != nil {
		return fmt.Errorf("subscribe to topic %q failed: %w", topic, err)
	}

	defer func() { _ = sub.Unsubscribe() }()

	<-ctx.Done()

	return nil
}

func (n *NATSMq) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

func connectNATS(url string) (*nats.Conn, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	return conn, nil
}
