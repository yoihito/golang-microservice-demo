package main

import (
	"context"
	"fmt"
	"gateway/pkg/handlers"
	"gateway/pkg/middlewares"
	"gateway/pkg/services"
	"gateway/pkg/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoDbUrl))
	if err != nil {
		log.Fatal(err)
	}

	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

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
		Auth:        authService,
		MongoClient: mongoClient,
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
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile("application.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	config := Config{
		Port:           viper.GetString("PORT"),
		MongoDbUrl:     viper.GetString("MONGO_DB_URL"),
		AuthServiceUrl: viper.GetString("AUTH_SERVICE_URL"),
	}
	return config, nil
}
