package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/gotasks/types"
	"github.com/google/uuid"
)

type DynamoDBProjectStore struct {
	client    *dynamodb.Client
	table     *string
	TaskStore *DynamoDBTaskStore
}

func NewDynamoDBProjectStore(client *dynamodb.Client, taskStore *DynamoDBTaskStore) *DynamoDBProjectStore {
	return &DynamoDBProjectStore{
		client:    client,
		table:     aws.String(projectColl),
		TaskStore: taskStore,
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
		return nil, ErrorNotFound
	}
	return project, nil
}

// TODO:
func (s *DynamoDBTaskStore) UpdateProjectTasks(ctx context.Context, id string, params Update) error {
	key, err := GetKey(id)
	if err != nil {
		return err
	}
	update := params.ToExpression()
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 s.table,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              dynamodbtypes.ReturnValueUpdatedNew,
	})
	if err != nil {
		return err
	}
	return nil
}
