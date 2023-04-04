package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDbUrl  string
	RabbitMqUrl string
	VideoQueue  string
	AudioQueue  string
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
		AudioQueue:  viper.GetString("AUDIO_QUEUE"),
	}
	return config, nil
}
