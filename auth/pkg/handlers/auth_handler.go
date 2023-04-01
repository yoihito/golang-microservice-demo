package handlers

import (
	"auth/pkg/utils"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type (
	JwtCustomClaims struct {
		Email string `json:"email"`
		jwt.RegisteredClaims
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	Handler struct {
		Db *sql.DB
	}
)

func (h *Handler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var passwordDigest string
	if err := h.Db.QueryRowContext(
		c.Request().Context(),
		"SELECT password_digest FROM users WHERE email = $1",
		req.Email).Scan(&passwordDigest); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": err.Error(),
		})
	}
	if !utils.MatchHashAndPassword(passwordDigest, req.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "email or password is invalid",
		})
	}

	t, err := createJWTToken(req.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func createJWTToken(email, secret string) (string, error) {
	claims := &JwtCustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (h *Handler) Validate(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	var userPresent bool
	if err := h.Db.QueryRowContext(
		c.Request().Context(),
		"SELECT 1 FROM users WHERE email = $1",
		claims.Email).Scan(&userPresent); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{})
}
