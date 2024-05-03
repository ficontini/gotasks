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

const (
	AwsProfileEnvName = "GOTASKS_AWS_PROFILE"
	dataTypeGSI       = "DataTypeGSI"
)

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
	//cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EnableAcceptEncodingGzip = true
	}), nil
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
	return map[string]dynamodbtypes.AttributeValue{dynamoIDField: id}, nil
}
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type DynamoDBQueryOptions struct {
	QueryInput *dynamodb.QueryInput
	Pagination *Pagination
}

func NewDynamoDBQueryOptions(queryInput *dynamodb.QueryInput, pagination *Pagination) *DynamoDBQueryOptions {
	return &DynamoDBQueryOptions{
		QueryInput: queryInput,
		Pagination: pagination,
	}
}

func PaginatedDynamoDBQuery(ctx context.Context, client *dynamodb.Client, opts *DynamoDBQueryOptions) ([]map[string]dynamodbtypes.AttributeValue, error) {
	var (
		collectiveResult []map[string]dynamodbtypes.AttributeValue
	)
	paginator := dynamodb.NewQueryPaginator(client, opts.QueryInput)
	for {
		if !paginator.HasMorePages() {
			break
		}
		singlePage, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		collectiveResult = append(collectiveResult, singlePage.Items...)
		if len(collectiveResult) >= (int(opts.Pagination.Page) * int(opts.Pagination.Limit)) {
			break
		}
	}
	return collectiveResult, nil
}
