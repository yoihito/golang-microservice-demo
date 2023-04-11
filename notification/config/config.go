package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	RabbitMqUrl  string
	AudioQueue   string
	FromEmail    string
	SmtpHost     string
	SmtpPort     string
	DownloadHost string
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
		RabbitMqUrl:  viper.GetString("RABBIT_MQ_URL"),
		AudioQueue:   viper.GetString("AUDIO_QUEUE"),
		FromEmail:    viper.GetString("FROM_EMAIL"),
		SmtpHost:     viper.GetString("SMTP_HOST"),
		SmtpPort:     viper.GetString("SMTP_PORT"),
		DownloadHost: viper.GetString("DOWNLOAD_HOST"),
	}
	return config, nil
}
