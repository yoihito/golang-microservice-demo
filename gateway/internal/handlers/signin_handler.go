package handlers

import (
	"gateway/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Signin(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return utils.JSONError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	token, err := h.Auth.Login(req.Email, req.Password)
	if err != nil {
		return utils.JSONError{Code: http.StatusUnauthorized, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
