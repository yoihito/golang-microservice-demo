package handlers

import (
	"gateway/internal/services"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type Handler struct {
	Auth           AuthService
	StorageService services.StorageService
	QueueService   services.QueueService
}
