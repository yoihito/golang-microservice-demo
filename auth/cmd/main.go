package main

import (
	"auth/pkg/controllers"
	"database/sql"
	"log"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func main() {
	connStr := "postgresql://authuser:mysecretpassword@localhost/auth?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	controller := &controllers.Controller{Db: db}

	e.POST("/login", controller.Login)
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(controllers.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	e.POST("/validate", controller.Validate, echojwt.WithConfig(config))
	e.Logger.Fatal(e.Start(":1323"))
}
