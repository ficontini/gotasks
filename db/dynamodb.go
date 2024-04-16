package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const AwsProfileEnvName = "AWS_PROFILE"

var PROFILE string

func NewDynamoDBStore() (*Store, error) {
	client, err := NewDynamoDBClient()
	if err != nil {
		return nil, err
	}
	task := NewDynamoDBTaskStore(client)
	return &Store{
		Auth:    NewDynamoDBAuthStore(client),
		User:    NewDynamoDBUserStore(client),
		Task:    task,
		Project: NewDynamoDBProjectStore(client, task),
	}, nil
}
func NewDynamoDBClient() (*dynamodb.Client, error) {
	if err := SetupDynamoDBConfigFromEnv(); err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(PROFILE))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func SetupDynamoDBConfigFromEnv() error {
	PROFILE = os.Getenv(AwsProfileEnvName)
	if PROFILE == "" {
		return fmt.Errorf("%s env variable not set", AwsProfileEnvName)
	}
	return nil
}
