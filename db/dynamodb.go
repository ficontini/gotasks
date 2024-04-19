package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const AwsProfileEnvName = "AWS_PROFILE"

var PROFILE string

func NewDynamoDBStore() (*Store, error) {
	client, err := NewDynamoDBClient()
	if err != nil {
		return nil, err
	}
	return &Store{
		Auth:    NewDynamoDBAuthStore(client),
		User:    NewDynamoDBUserStore(client),
		Task:    NewDynamoDBTaskStore(client),
		Project: NewDynamoDBProjectStore(client),
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

func GetKey(idStr string) (map[string]dynamodbtypes.AttributeValue, error) {
	id, err := attributevalue.Marshal(idStr)
	if err != nil {
		return nil, err
	}
	return map[string]dynamodbtypes.AttributeValue{"ID": id}, nil
}
