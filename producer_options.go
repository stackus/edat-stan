package edatstan

import (
	"github.com/stackus/edat/log"
)

type ProducerOption func(*Producer)

func WithProducerSerializer(serializer Serializer) ProducerOption {
	return func(producer *Producer) {
		producer.serializer = serializer
	}
}

func WithProducerLogger(logger log.Logger) ProducerOption {
	return func(producer *Producer) {
		producer.logger = logger
	}
}
