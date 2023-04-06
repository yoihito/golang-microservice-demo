package main

import (
	"fmt"
	"log"
	"time"

	"converter/internal/operations"
	"converter/internal/services"

	"converter/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := services.NewGridFSService(config.MongoDbUrl)
	if err != nil {
		log.Fatal(err)
	}

	consumerQueue := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.VideoQueue}})
	publisherQueue := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.AudioQueue}})

	for {
		if consumerQueue.IsReady() && publisherQueue.IsReady() {
			break
		}
		<-time.After(10 * time.Second)
	}

	msgs, err := consumerQueue.Consume(config.VideoQueue, false)
	if err != nil {
		log.Fatal(err)
	}

	processor := operations.Processor{
		Storage: storage,
		Queue:   publisherQueue,
		Config:  config,
	}
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body())
			if err := processor.ProcessMessage(msg); err != nil {
				msg.Nack()
				continue
			}
			msg.Ack()
		}
	}()
	<-forever
}
