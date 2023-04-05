package services

import (
	"context"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type StorageService interface {
	UploadFromStream(database, filename string, reader io.Reader) (string, error)
}

type GridFSService struct {
	mongoClient *mongo.Client
}

func NewGridFSService(serviceUrl string) (*GridFSService, error) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(serviceUrl).SetTimeout(5*time.Second))
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

func (s *GridFSService) UploadFromStream(database, filename string, src io.Reader) (string, error) {
	db := s.mongoClient.Database(database)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return "", err
	}

	objectId, err := bucket.UploadFromStream(filename, src)

	if err != nil {
		return "", err
	}

	return objectId.Hex(), nil
}

func (s *GridFSService) Close() {
	s.mongoClient.Disconnect(context.TODO())
}
