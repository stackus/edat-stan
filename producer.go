package edatstan

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

type Producer struct {
	conn       stan.Conn
	serializer Serializer
	logger     log.Logger
}

var _ msg.Producer = (*Producer)(nil)

func NewProducer(conn stan.Conn, options ...ProducerOption) *Producer {
	p := &Producer{
		conn:       conn,
		serializer: DefaultSerializer,
		logger:     log.DefaultLogger,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *Producer) Send(ctx context.Context, channel string, message msg.Message) error {
	logger := p.logger.Sub(
		log.String("Channel", channel),
	)

	data, err := p.serializer.Serialize(message)
	if err != nil {
		logger.Error("failed to marshal message", log.Error(err))
		return fmt.Errorf("message could not be marshalled")
	}

	if err = p.conn.Publish(channel, data); err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close(context.Context) error {
	p.logger.Trace("closing message destination")
	err := p.conn.Close()
	if err != nil {
		p.logger.Error("error closing message destination", log.Error(err))
	}
	return err
}
