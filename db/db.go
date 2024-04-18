package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoEndpoint      = "MONGO_DB_URI"
	MongoDBNameEnvName = "MONGO_DB_NAME"
)

var (
	DBNAME string
	DBURI  string
)

type Store struct {
	Auth    AuthStore
	User    UserStore
	Task    TaskStore
	Project ProjectStore
}

func NewStore() (*Store, error) {
	client, err := NewMongoClient()
	if err != nil {
		return nil, err
	}
	dynamoClient, err := NewDynamoDBClient()
	if err != nil {
		return nil, err
	}
	taskStore := NewMongoTaskStore(client)
	return &Store{
		Auth:    NewDynamoDBAuthStore(dynamoClient),
		User:    NewMongoUserStore(client),
		Task:    taskStore,
		Project: NewMongoProjectStore(client, taskStore),
	}, nil
}
func NewMongoClient() (*mongo.Client, error) {
	if err := SetupMongoConfigFromEnv(); err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SetupMongoConfigFromEnv() error {
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
