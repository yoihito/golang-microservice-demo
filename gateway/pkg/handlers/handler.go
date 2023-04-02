package handlers

import (
	"gateway/pkg/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type Handler struct {
	Auth         AuthService
	MongoClient  *mongo.Client
	QueueService services.QueueService
}
