package handlers

import (
	"auth/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		h.Logger.Warn().Err(err).Msg("failed to bind request")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		h.Logger.Warn().Err(err).Msg("failed to validate request")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, err := h.Repo.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("user not found")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "email or password is invalid",
		})
	}
	if !utils.MatchHashAndPassword(user.PasswordDigest, req.Password) {
		h.Logger.Warn().Err(err).Msg("passwords doesn't match")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "email or password is invalid",
		})
	}

	if t, err := h.TokenManager.CreateToken(user.Email); err != nil {
		h.Logger.Error().Err(err).Msg("failed to create token")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to create token",
		})
	} else {
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	}
}

func (h *Handler) Validate(c echo.Context) error {
	email, ok := h.TokenManager.GetEmailFromToken(c)
	if !ok {
		h.Logger.Warn().Msg("token is invalid")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "token is invalid",
		})
	}

	if user, err := h.Repo.GetByEmail(c.Request().Context(), email); err != nil {
		h.Logger.Warn().Err(err).Msg("user not found")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "token is invalid",
		})
	} else {
		return c.JSON(http.StatusOK, echo.Map{
			"email": user.Email,
		})
	}
}
