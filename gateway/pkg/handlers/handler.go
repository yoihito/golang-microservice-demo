package handlers

import "go.mongodb.org/mongo-driver/mongo"

type AuthService interface {
	Login(email, password string) (string, error)
	Validate(token string) error
}

type Handler struct {
	Auth        AuthService
	MongoClient *mongo.Client
}
