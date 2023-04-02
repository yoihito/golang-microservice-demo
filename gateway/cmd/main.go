package main

import (
	"fmt"
	"gateway/pkg/handlers"
	"gateway/pkg/middlewares"
	"gateway/pkg/services"
	"gateway/pkg/utils"
	"io/fs"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	storageService, err := services.NewGridFSService(config.MongoDbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer storageService.Close()

	queueService, err := services.NewRabbitMqService(config.RabbitMqUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer queueService.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if jsonError, ok := err.(utils.JSONError); ok {
			if err := c.JSON(jsonError.Code, jsonError); err != nil {
				e.Logger.Error(err)
			}
			return
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	authService := services.NewAuthService(
		config.AuthServiceUrl,
	)
	handler := &handlers.Handler{
		Auth:           authService,
		StorageService: storageService,
		QueueService:   queueService,
	}
	e.POST("/signin", handler.Signin)

	restricted := e.Group("")
	restricted.Use(middlewares.TokenVerification(authService))
	restricted.POST("/upload", handler.Upload)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}

type Config struct {
	MongoDbUrl     string
	Port           string
	AuthServiceUrl string
	RabbitMqUrl    string
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
		Port:           viper.GetString("PORT"),
		MongoDbUrl:     viper.GetString("MONGO_DB_URL"),
		AuthServiceUrl: viper.GetString("AUTH_SERVICE_URL"),
		RabbitMqUrl:    viper.GetString("RABBITMQ_URL"),
	}
	return config, nil
}
