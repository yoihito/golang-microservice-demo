package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl string
	SecretKey   string
	Port        string
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
		DatabaseUrl: viper.GetString("DATABASE_URL"),
		SecretKey:   viper.GetString("JWT_SECRET"),
		Port:        viper.GetString("PORT"),
	}
	return config, nil
}
