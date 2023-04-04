package main

import (
	"fmt"
	"log"

	"converter/internal/operations"
	"converter/internal/services"

	"converter/internal/config"
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

	consumerQueue, err := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.VideoQueue}})
	publisherQueue, err := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.AudioQueue}})

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

	fmt.Println("Waiting for messages...")
	<-forever
}
