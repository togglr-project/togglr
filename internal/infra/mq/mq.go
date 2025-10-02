package mq

import (
	"context"
)

type MQ interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(
		ctx context.Context,
		topic string,
		processFn func(ctx context.Context, data []byte) error,
	) error
	Request(
		ctx context.Context,
		topic string,
		data []byte,
	) ([]byte, error)
	HandleRequest(
		ctx context.Context,
		topic string,
		handlerFn func(ctx context.Context, data []byte) (response []byte, err error),
	) error
}
