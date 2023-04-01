package middlewares

import (
	"gateway/pkg/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthService interface {
	Validate(token string) error
}

func TokenVerification(auth AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			parts := strings.Split(token, " ")
			if len(parts) < 2 {
				return utils.JSONError{
					Code:    401,
					Message: "Invalid token",
				}
			}
			err := auth.Validate(parts[1])
			if err != nil {
				return utils.JSONError{
					Code:    401,
					Message: err.Error(),
				}
			}
			return next(c)
		}
	}
}
