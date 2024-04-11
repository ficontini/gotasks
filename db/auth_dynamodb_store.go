package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/gotasks/types"
)

const authTable = "auths"

type DynamoDBAuthStore struct {
	client *dynamodb.Client
}

func NewDynamoDBAuthStore(client *dynamodb.Client) *DynamoDBAuthStore {
	return &DynamoDBAuthStore{
		client: client,
	}
}

func (s *DynamoDBAuthStore) Insert(ctx context.Context, auth *types.Auth) (*types.Auth, error) {
	item, err := attributevalue.MarshalMap(auth)
	if err != nil {
		return nil, err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(authTable), Item: item,
	})
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	return auth, nil
}

func (s *DynamoDBAuthStore) Get(ctx context.Context, params *types.AuthFilter) (*types.Auth, error) {
	key, err := attributevalue.MarshalMap(params)
	if err != nil {
		return nil, err
	}
	res, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(authTable),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, ErrorNotFound
	}
	var auth *types.Auth
	if err := attributevalue.UnmarshalMap(res.Item, &auth); err != nil {
		return nil, err
	}
	fmt.Println("auth", auth)
	return auth, nil
}

func (s *DynamoDBAuthStore) Delete(ctx context.Context, params *types.AuthFilter) error {
	key, err := attributevalue.MarshalMap(params)
	if err != nil {
		return err
	}
	_, err = s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(authTable),
	})
	if err != nil {
		return err
	}
	return nil
}
