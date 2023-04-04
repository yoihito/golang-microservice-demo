package main

import (
	"fmt"
	"io/fs"
	"log"

	"converter/internal/operations"
	"converter/internal/services"

	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := services.NewGridFSService(config.MongoDbUrl)
	if err != nil {
		log.Fatal(err)
	}

	consumerQueue, err := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.VideoQueue}})
	publisherQueue, err := services.NewRabbitMqService(config.RabbitMqUrl, []services.RabbitMqQueue{{config.AudioQueue()}})

	msgs, err := consumerQueue.Consume(config.VideoQueue, false)
	if err != nil {
		log.Fatal(err)
	}

	processor := operations.Processor{
		Storage: storage,
		Queue:   publisherQueue,
		Config:  &config,
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

type Config struct {
	MongoDbUrl  string
	RabbitMqUrl string
	VideoQueue  string
	audioQueue  string
}

func (c *Config) AudioQueue() string {
	return c.audioQueue
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile("application.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *fs.PathError:
		default:
			return Config{}, nil
		}
	}

	config := Config{
		MongoDbUrl:  viper.GetString("MONGO_DB_URL"),
		RabbitMqUrl: viper.GetString("RABBIT_MQ_URL"),
		VideoQueue:  viper.GetString("VIDEO_QUEUE"),
		audioQueue:  viper.GetString("AUDIO_QUEUE"),
	}
	return config, nil
}
