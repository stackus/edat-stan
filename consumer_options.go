package edatstan

import (
	"time"

	"github.com/nats-io/stan.go"

	"github.com/stackus/edat/log"
)

type ConsumerOption func(*Consumer)

func WithConsumerActWait(ackWait time.Duration) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.ackWait = ackWait
		consumer.subOptions = append(consumer.subOptions, stan.AckWait(ackWait))
	}
}

func WithConsumerSubscriptionOptions(option ...stan.SubscriptionOption) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.subOptions = append(consumer.subOptions, option...)
	}
}

func WithConsumerSerializer(serializer Serializer) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.serializer = serializer
	}
}

func WithConsumerLogger(logger log.Logger) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.logger = logger
	}
}
