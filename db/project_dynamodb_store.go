package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/gotasks/types"
	"github.com/google/uuid"
)

type DynamoDBProjectStore struct {
	client    *dynamodb.Client
	table     *string
	taskStore *DynamoDBTaskStore
}

func NewDynamoDBProjectStore(client *dynamodb.Client, taskStore *DynamoDBTaskStore) *DynamoDBProjectStore {
	return &DynamoDBProjectStore{
		client:    client,
		table:     aws.String(projectColl),
		taskStore: taskStore,
	}
}
func (s *DynamoDBProjectStore) InsertProject(ctx context.Context, project *types.Project) (*types.Project, error) {
	project.ID = uuid.New().String()
	item, err := attributevalue.MarshalMap(project)
	if err != nil {
		return nil, err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: s.table, Item: item,
	})
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *DynamoDBProjectStore) GetProjectByID(ctx context.Context, id string) (*types.Project, error) {
	key, err := GetKey(id)
	if err != nil {
		return nil, err
	}
	res, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: s.table,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, ErrorNotFound
	}
	var project *types.Project
	if err := attributevalue.UnmarshalMap(res.Item, &project); err != nil {
		return nil, err
	}
	return project, nil
}
func (s *DynamoDBProjectStore) TransactWriteItems(ctx context.Context, actions []UpdateAction) error {
	operations := make([]dynamodbtypes.TransactWriteItem, 0, len(actions))
	for _, action := range actions {
		operation, err := action.get()
		if err != nil {
			return err
		}
		update, ok := operation.(*dynamodbtypes.Update)
		if !ok {
			return ErrInvalidOperationType
		}
		writeItem := dynamodbtypes.TransactWriteItem{
			Update: update,
		}
		operations = append(operations, writeItem)

	}
	_, err := s.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: operations,
	})
	if err != nil {
		return err
	}
	return nil
}

var ErrInvalidOperationType = errors.New("invalid operation")
