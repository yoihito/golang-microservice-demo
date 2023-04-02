package services

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type StorageService interface {
	UploadFromStream(filename string, reader io.Reader) (string, error)
}

type GridFSService struct {
	mongoClient *mongo.Client
}

func NewGridFSService(serviceUrl string) (*GridFSService, error) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(serviceUrl))
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &GridFSService{
		mongoClient: mongoClient,
	}, nil
}

func (s *GridFSService) UploadFromStream(filename string, src io.Reader) (string, error) {
	db := s.mongoClient.Database("videos")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return "", err
	}

	objectId, err := bucket.UploadFromStream(filename, src)

	if err != nil {
		return "", err
	}

	return objectId.String(), nil
}

func (s *GridFSService) Close() {
	s.mongoClient.Disconnect(context.TODO())
}
