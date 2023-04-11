package utils

import (
	"github.com/labstack/echo/v4"
)

func NewHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if jsonError, ok := err.(JSONError); ok {
			if err := c.JSON(jsonError.Code, jsonError); err != nil {
				e.Logger.Error(err)
			}
			return
		}
		e.DefaultHTTPErrorHandler(err, c)
	}
}
