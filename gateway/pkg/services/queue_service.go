package services

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService interface {
	Publish(ctx context.Context, topic string, message []byte) error
}

type RabbitMqService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMqService(serviceUrl string) (*RabbitMqService, error) {
	conn, err := amqp.Dial(serviceUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMqService{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMqService) Publish(ctx context.Context, topic string, data []byte) error {
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
