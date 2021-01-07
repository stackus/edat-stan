package edatstan

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

var DefaultAckWait = time.Second * 30

type Consumer struct {
	conn       stan.Conn
	queue      string
	ackWait    time.Duration
	serializer Serializer
	subOptions []stan.SubscriptionOption
	listenerWg sync.WaitGroup
	logger     log.Logger
}

var _ msg.Consumer = (*Consumer)(nil)

func NewConsumer(conn stan.Conn, groupID string, options ...ConsumerOption) *Consumer {
	c := &Consumer{
		conn:  conn,
		queue: groupID,
		subOptions: []stan.SubscriptionOption{
			stan.SetManualAckMode(),
			stan.AckWait(DefaultAckWait),
			stan.DurableName("Durable"),
		},
		ackWait:    DefaultAckWait,
		serializer: DefaultSerializer,
		logger:     log.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Consumer) Listen(ctx context.Context, channel string, subscription msg.ReceiveMessageFunc) error {
	logger := c.logger.Sub(log.String("Channel", channel))

	defer logger.Trace("stopped listening")

	_, err := c.conn.QueueSubscribe(channel, c.queue, c.consumeMessages(ctx, subscription), c.subOptions...)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (c *Consumer) Close(context.Context) error {
	c.logger.Trace("closing message source")
	err := c.conn.Close()
	if err != nil {
		c.logger.Error("error closing message source", log.Error(err))
	}

	return err
}

func (c *Consumer) consumeMessages(ctx context.Context, receiver func(context.Context, msg.Message) error) func(*stan.Msg) {
	return func(stanMsg *stan.Msg) {
		var err error

		select {
		case <-ctx.Done():
			c.logger.Trace("listener has closed; message processing stopped")
			return
		default:
		}

		var message msg.Message
		message, err = c.serializer.Deserialize(stanMsg)
		if err != nil {
			c.logger.Error("message failed to unmarshal", log.Error(err))
			return
		}

		wCtx, cancel := context.WithTimeout(ctx, c.ackWait)
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- receiver(wCtx, message)
		}()

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := stanMsg.Ack(); ackErr != nil {
					c.logger.Error("error acknowledging message", log.Error(err))
				}
			}
		case <-wCtx.Done():
			c.logger.Trace("listener has closed; in-progress message processing terminated")
		}
	}
}
