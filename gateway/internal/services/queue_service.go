package services

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService interface {
	Publish(ctx context.Context, topic string, message []byte) error
}

type RabbitMqService struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	queues        []RabbitMqQueue
	notifyConfirm chan amqp.Confirmation
	isReady       bool
}

type RabbitMqQueue struct {
	Name string
}

func NewRabbitMqService(serviceUrl string, queues []RabbitMqQueue) (*RabbitMqService, error) {
	s := &RabbitMqService{queues: queues}
	if err := s.connect(serviceUrl); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *RabbitMqService) connect(serviceUrl string) error {
	c := make(chan *amqp.Error)
	go func() {
		<-c
		r.connect(serviceUrl)
	}()
	conn, err := amqp.Dial(serviceUrl)
	if err != nil {
		return err
	}
	r.connection = conn
	if err := r.init(); err != nil {
		return err
	}
	conn.NotifyClose(c)
	return nil
}

func (r *RabbitMqService) init() error {
	c := make(chan *amqp.Error)
	go func() {
		<-c
		r.init()
	}()
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
	r.notifyConfirm = make(chan amqp.Confirmation)
	ch.NotifyClose(c)
	ch.NotifyPublish(r.notifyConfirm)
	r.isReady = true
	return nil
}

func (r *RabbitMqService) Publish(ctx context.Context, topic string, data []byte) error {
	if !r.isReady {
		return errors.New("failed to push: not connected")
	}
	for {
		err := r.UnsafePublish(ctx, topic, data)
		if err != nil {
			<-time.After(5 * time.Second)
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

func (r *RabbitMqService) Close() {
	r.channel.Close()
	r.connection.Close()
}
