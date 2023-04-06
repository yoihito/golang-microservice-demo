package app

import (
	"auth/config"
	"auth/internal/handlers"
	"auth/internal/infrustructure"
	"auth/internal/repositories"
	"auth/internal/utils"
	"fmt"
	"log"

	"github.com/go-playground/validator"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(config config.Config) {
	datastore, err := infrustructure.NewDatastore(config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Validator = utils.NewCustomValidator(validator.New())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	tokenManager := infrustructure.NewJWTTokenManager(config.SecretKey)

	handler := &handlers.Handler{
		Repo:         repositories.NewUser(datastore),
		TokenManager: tokenManager,
	}

	e.POST("/login", handler.Login)
	e.POST("/validate", handler.Validate, echojwt.WithConfig(tokenManager.JWTMiddlewareConfig()))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
