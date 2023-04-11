package handlers

import (
	"auth/internal/entities"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UserRepo interface {
	GetByEmail(context.Context, string) (*entities.User, error)
}

type TokenManager interface {
	CreateToken(email string) (string, error)
	GetEmailFromToken(c echo.Context) (string, bool)
}

type Handler struct {
	Repo         UserRepo
	TokenManager TokenManager
	Logger       *zerolog.Logger
}
