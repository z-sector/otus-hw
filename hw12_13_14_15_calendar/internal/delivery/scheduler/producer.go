package scheduler

import (
	"context"
	"encoding/json"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/amqp"
)

type ProducerRMQ struct {
	client *amqp.ClientAMQP
}

func (p *ProducerRMQ) Publish(ctx context.Context, notif internal.EventNotification) error {
	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	return p.client.Push(ctx, data, notif.EventID.String())
}

func NewProducerRMQ(client *amqp.ClientAMQP) *ProducerRMQ {
	return &ProducerRMQ{client: client}
}
