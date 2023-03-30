package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	e := echo.New()
	e.POST("/login", func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")
		if email == "" || password == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "email or password is required",
			})
		}
		if email != "test@example.com" || password != "test" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "email or password is invalid",
			})
		}
		claims := &jwtCustomClaims{
			email,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
