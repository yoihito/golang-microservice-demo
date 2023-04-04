package services

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService interface {
	Publish(ctx context.Context, topic string, message any) error
	Consume(topic string, autoAck bool) (<-chan Delivery, error)
}

type Delivery interface {
	Body() []byte
	Ack() error
	Nack() error
}

type RabbitMqDelivery struct {
	msg amqp.Delivery
}

func (d *RabbitMqDelivery) Body() []byte {
	return d.msg.Body
}

func (d *RabbitMqDelivery) Ack() error {
	return d.msg.Ack(false)
}

func (d *RabbitMqDelivery) Nack() error {
	return d.msg.Nack(false, true)
}

type RabbitMqService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

type RabbitMqQueue struct {
	Name string
}

func NewRabbitMqService(serviceUrl string, queues []RabbitMqQueue) (*RabbitMqService, error) {
	conn, err := amqp.Dial(serviceUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, queue := range queues {
		_, err = ch.QueueDeclare(queue.Name, false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	return &RabbitMqService{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMqService) Publish(ctx context.Context, topic string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx, "", topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
	})
}

func (r *RabbitMqService) Consume(topic string, autoAck bool) (<-chan Delivery, error) {
	msgs, err := r.channel.Consume("videos", "", autoAck, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	channel := make(chan Delivery)
	go func() {
		for msg := range msgs {
			delivery := &RabbitMqDelivery{msg}
			channel <- delivery
		}
	}()
	return channel, nil
}

func (r *RabbitMqService) Close() {
	r.channel.Close()
	r.connection.Close()
}
