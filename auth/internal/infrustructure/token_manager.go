package infrustructure

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTTokenManager struct {
	secretKey string
}

func NewJWTTokenManager(secretKey string) *JWTTokenManager {
	return &JWTTokenManager{
		secretKey: secretKey,
	}
}

func (m *JWTTokenManager) CreateToken(email string) (string, error) {
	claims := &JwtCustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (m *JWTTokenManager) JWTMiddlewareConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(m.secretKey),
	}
}

func (m *JWTTokenManager) GetEmailFromToken(c echo.Context) (string, bool) {
	if user, ok := c.Get("user").(*jwt.Token); ok {
		if claims, ok := user.Claims.(*JwtCustomClaims); ok {
			return claims.Email, true
		}
	}

	return "", false
}
