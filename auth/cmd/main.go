package main

import (
	"auth/internal/handlers"
	"auth/internal/infrustructure"
	"auth/internal/repositories"
	"auth/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	datastore, err := infrustructure.NewDatastore(connStr)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Validator = utils.NewCustomValidator(validator.New())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	tokenManager := infrustructure.NewJWTTokenManager(os.Getenv("JWT_SECRET"))

	handler := &handlers.Handler{
		Repo:         repositories.NewUser(datastore),
		TokenManager: tokenManager,
	}

	e.POST("/login", handler.Login)
	e.POST("/validate", handler.Validate, echojwt.WithConfig(tokenManager.JWTMiddlewareConfig()))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
