package main

import (
	"fmt"
	"gateway/config"
	"gateway/internal/handlers"
	"gateway/internal/middlewares"
	"gateway/internal/services"
	"gateway/internal/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	storageService, err := services.NewGridFSService(config.MongoDbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer storageService.Close()

	queues := []services.RabbitMqQueue{}
	for _, queueName := range config.Queues {
		queues = append(queues, services.RabbitMqQueue{Name: queueName})
	}

	queueService := services.NewRabbitMqService(config.RabbitMqUrl, queues)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
