package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoEndpoint      = "MONGO_DB_URI"
	MongoDBNameEnvName = "MONGO_DB_NAME"
	AwsProfileEnvName  = "AWS_PROFILE"
)

var (
	DBNAME        string
	DBURI         string
	PROFILE       string
	ErrorNotFound = errors.New("resource not found")
	ErrInvalidID  = errors.New("invalid ID")
)

type Store struct {
	Auth    AuthStore
	User    UserStore
	Task    TaskStore
	Project ProjectStore
}

func NewStore() (*Store, error) {
	if err := SetupDBConfigFromEnv(); err != nil {
		return nil, err
	}
	client, err := newMongoClient()
	if err != nil {
		return nil, err
	}
	dynamoClient, err := newDynamoDBClient()
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

func SetupDBConfigFromEnv() error {
	DBURI = os.Getenv(MongoEndpoint)
	if DBURI == "" {
		return fmt.Errorf("%s env variable not set", MongoEndpoint)
	}
	DBNAME = os.Getenv(MongoDBNameEnvName)
	if DBNAME == "" {
		return fmt.Errorf("%s env variable not set", MongoDBNameEnvName)
	}
	PROFILE = os.Getenv(AwsProfileEnvName)
	if PROFILE == "" {
		return fmt.Errorf("%s env variable not set", AwsProfileEnvName)
	}
	return nil
}
func newMongoClient() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		return nil, err
	}
	return client, nil
}
func newDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(PROFILE))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}
