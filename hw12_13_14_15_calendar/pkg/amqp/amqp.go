package amqp

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

const (
	defaultConnAttempts = 20
	reconnectDelay      = 5 * time.Second
	reInitDelay         = 2 * time.Second
	resendDelay         = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

type ClientAMQP struct {
	queueName       string
	prefetchCount   int
	logger          logger.AppLog
	connection      *amqp.Connection
	Channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

func NewClientAMQP(log logger.AppLog, addr, queueName string) *ClientAMQP {
	client := ClientAMQP{
		logger:    log,
		queueName: queueName,
		done:      make(chan bool),
	}
	go client.handleReconnect(addr)
	return &client
}

func (client *ClientAMQP) WithPrefetchCount(count int) *ClientAMQP {
	client.prefetchCount = count
	return client
}

func (client *ClientAMQP) handleReconnect(addr string) {
	for {
		client.isReady = false
		client.logger.Info("Attempting to connect")

		conn, err := client.connect(addr)
		if err != nil {
			client.logger.Info("Failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

func (client *ClientAMQP) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.logger.Info("Connected!")
	return conn, nil
}

func (client *ClientAMQP) handleReInit(conn *amqp.Connection) bool {
	for {
		client.isReady = false

		err := client.init(conn)
		if err != nil {
			client.logger.Info("Failed to initialize channel. Retrying...")

			select {
			case <-client.done:
				return true
			case <-client.notifyConnClose:
				client.logger.Info("Connection closed. Reconnecting...")
				return false
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.logger.Info("Connection closed. Reconnecting...")
			return false
		case <-client.notifyChanClose:
			client.logger.Info("Channel closed. Re-running init...")
		}
	}
}

func (client *ClientAMQP) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		client.queueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-message-deduplication": true},
	)

	if err != nil {
		return err
	}

	client.changeChannel(ch)
	client.isReady = true
	client.logger.Info("Setup!")

	return nil
}

func (client *ClientAMQP) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

func (client *ClientAMQP) changeChannel(channel *amqp.Channel) {
	client.Channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.Channel.NotifyClose(client.notifyChanClose)
	client.Channel.NotifyPublish(client.notifyConfirm)
}

func (client *ClientAMQP) Push(ctx context.Context, data []byte, deduplicationKey string) error {
	if !client.isReady {
		return errors.New("failed to push: not connected")
	}
	for {
		err := client.UnsafePush(ctx, data, deduplicationKey)
		if err != nil {
			client.logger.Info("Push failed. Retrying...")
			select {
			case <-client.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		confirm := <-client.notifyConfirm
		if confirm.Ack {
			client.logger.Info(fmt.Sprintf("Push confirmed [%d]!", confirm.DeliveryTag))
			return nil
		}
	}
}

func (client *ClientAMQP) UnsafePush(ctx context.Context, data []byte, deduplicationKey string) error {
	if !client.isReady {
		return errNotConnected
	}

	return client.Channel.PublishWithContext(
		ctx,
		"",
		client.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Headers:      amqp.Table{"x-deduplication-header": deduplicationKey},
			ContentType:  "application/json",
			Body:         data,
		},
	)
}

func (client *ClientAMQP) Consume() (<-chan amqp.Delivery, error) {
	connAttempts := defaultConnAttempts
	for connAttempts > 0 {
		if client.isReady {
			break
		}
		client.logger.Info(fmt.Sprintf("Rabbit is trying to connect, attempts left: %d", connAttempts))
		time.Sleep(time.Second)
		connAttempts--
	}

	if !client.isReady {
		return nil, errNotConnected
	}

	if err := client.Channel.Qos(
		client.prefetchCount,
		0,
		false,
	); err != nil {
		return nil, err
	}

	return client.Channel.Consume(
		client.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (client *ClientAMQP) Close() error {
	if !client.isReady {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.Channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	return nil
}
