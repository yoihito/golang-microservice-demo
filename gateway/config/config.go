package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDbUrl     string
	Port           string
	AuthServiceUrl string
	RabbitMqUrl    string
	Queues         []string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile("./config/application.yaml")
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
		Port:           viper.GetString("PORT"),
		MongoDbUrl:     viper.GetString("MONGO_DB_URL"),
		AuthServiceUrl: viper.GetString("AUTH_SERVICE_URL"),
		RabbitMqUrl:    viper.GetString("RABBIT_MQ_URL"),
		Queues:         viper.GetStringSlice("QUEUES"),
	}
	return config, nil
}
