package services

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	errNotConnected  = errors.New("not connected to RabbitMQ server")
	errAlreadyClosed = errors.New("already closed")
	errShutdown      = errors.New("shutting down")
)

type QueueService interface {
	Publish(ctx context.Context, topic string, message []byte) error
}

type RabbitMqService struct {
	connection      *amqp.Connection
	channel         *amqp.Channel
	queues          []RabbitMqQueue
	done            chan struct{}
	notifyConnClose chan *amqp.Error
	notifyChClose   chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

type RabbitMqQueue struct {
	Name string
}

func NewRabbitMqService(serviceUrl string, queues []RabbitMqQueue) *RabbitMqService {
	s := &RabbitMqService{queues: queues}
	go s.Reconnect(serviceUrl)
	return s
}

func (r *RabbitMqService) Reconnect(serviceUrl string) {
	for {
		r.isReady = false
		err := r.connect(serviceUrl)
		if err != nil {
			select {
			case <-r.done:
				return
			case <-time.After(5 * time.Second):
			}

			continue
		}
		if done := r.reInit(); done {
			break
		}
	}
}

func (r *RabbitMqService) connect(serviceUrl string) error {
	conn, err := amqp.Dial(serviceUrl)
	if err != nil {
		return err
	}
	r.connection = conn
	r.notifyConnClose = make(chan *amqp.Error, 1)
	r.connection.NotifyClose(r.notifyConnClose)
	return nil
}

func (r *RabbitMqService) reInit() bool {
	for {
		r.isReady = false
		err := r.init()
		if err != nil {
			select {
			case <-r.done:
				return true
			case <-time.After(5 * time.Second):
			}
			continue
		}
		select {
		case <-r.done:
			return true
		case <-r.notifyConnClose:
			return false
		case <-r.notifyChClose:
		}
	}
}

func (r *RabbitMqService) init() error {
	ch, err := r.connection.Channel()
	if err != nil {
		return err
	}
	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	for _, queue := range r.queues {
		_, err = ch.QueueDeclare(queue.Name, false, false, false, false, nil)
		if err != nil {
			return err
		}
	}
	r.channel = ch
	r.notifyChClose = make(chan *amqp.Error, 1)
	r.notifyConfirm = make(chan amqp.Confirmation, 1)
	ch.NotifyClose(r.notifyChClose)
	ch.NotifyPublish(r.notifyConfirm)
	r.isReady = true
	return nil
}

func (r *RabbitMqService) Publish(ctx context.Context, topic string, data []byte) error {
	if !r.isReady {
		return errNotConnected
	}
	for {
		err := r.UnsafePublish(ctx, topic, data)
		if err != nil {
			select {
			case <-r.done:
				return errShutdown
			case <-time.After(5 * time.Second):
			}
			continue
		}
		confirm := <-r.notifyConfirm
		if confirm.Ack {
			return nil
		}
	}
}

func (r *RabbitMqService) UnsafePublish(ctx context.Context, topic string, data []byte) error {
	return r.channel.PublishWithContext(ctx, "", topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
	})
}

func (r *RabbitMqService) Close() error {
	if !r.isReady {
		return errAlreadyClosed
	}
	close(r.done)
	r.channel.Close()
	r.connection.Close()
	r.isReady = false
	return nil
}
