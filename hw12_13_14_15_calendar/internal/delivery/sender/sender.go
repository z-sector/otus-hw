package sender

import (
	"context"
	"fmt"

	origamqp "github.com/rabbitmq/amqp091-go"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/amqp"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Subscriber interface {
	Consume(ctx context.Context, data []byte) error
}

type Sender struct {
	log    logger.AppLog
	client *amqp.ClientAMQP
}

func NewSender(log logger.AppLog, client *amqp.ClientAMQP) *Sender {
	return &Sender{log: log, client: client}
}

func (s *Sender) Run(ctx context.Context, subscriber Subscriber) error {
	deliveries, err := s.client.Consume()
	if err != nil {
		return err
	}
	chClosedCh := make(chan *origamqp.Error, 1)
	s.client.Channel.NotifyClose(chClosedCh)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case amqErr := <-chClosedCh:
			s.log.Info(fmt.Sprintf("AMQP Channel closed due to: %s", amqErr))

			deliveries, err = s.client.Consume()
			if err != nil {
				s.log.Info("Error trying to consume, will try again")
				continue
			}

			chClosedCh = make(chan *origamqp.Error, 1)
			s.client.Channel.NotifyClose(chClosedCh)

		case delivery := <-deliveries:
			s.log.Info(fmt.Sprintf("Received message: %s", delivery.Body))
			if err := subscriber.Consume(ctx, delivery.Body); err != nil {
				s.log.Error("Error handle body", err)
			}
			if err := delivery.Ack(false); err != nil {
				s.log.Error("Error acknowledging message", err)
			}
		}
	}
}

func (s *Sender) Close() error {
	return s.client.Close()
}
