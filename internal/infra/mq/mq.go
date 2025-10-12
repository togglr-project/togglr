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
	SubscribeBatch(
		ctx context.Context,
		streamName string,
		batchSize int,
		processFn func(ctx context.Context, messages [][]byte) error,
	) error
}
