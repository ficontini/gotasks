package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DEFAULT_PAGE       = 1
	DEFAULT_LIMIT      = 10
	MongoEndpoint      = "MONGO_DB_URI"
	MongoDBNameEnvName = "MONGO_DB_NAME"
)

var (
	DBNAME        string
	DBURI         string
	ErrorNotFound = errors.New("resource not found")
)

type Store struct {
	User    UserStore
	Task    TaskStore
	Project ProjectStore
}

func NewMongoStore() (*Store, error) {
	if err := SetupMongoDBConfigFromEnv(); err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		return nil, err
	}
	taskStore := NewMongoTaskStore(client)
	return &Store{
		User:    NewMongoUserStore(client),
		Task:    taskStore,
		Project: NewMongoProjectStore(client, taskStore),
	}, nil
}

func SetupMongoDBConfigFromEnv() error {
	DBURI = os.Getenv(MongoEndpoint)
	if DBURI == "" {
		return fmt.Errorf("%s env variable not set", MongoEndpoint)
	}
	DBNAME = os.Getenv(MongoDBNameEnvName)
	if DBNAME == "" {
		return fmt.Errorf("%s env variable not set", MongoDBNameEnvName)
	}
	return nil
}

type Deleter interface {
	Delete(context.Context, string) error
}
