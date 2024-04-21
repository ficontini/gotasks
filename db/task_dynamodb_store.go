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

const ReturnAllOld = "ALL_OLD"

type DynamoDBTaskStore struct {
	client *dynamodb.Client
	table  *string
}

func NewDynamoDBTaskStore(client *dynamodb.Client) *DynamoDBTaskStore {
	return &DynamoDBTaskStore{
		client: client,
		table:  aws.String(taskColl),
	}
}

func (s *DynamoDBTaskStore) InsertTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	task.ID = uuid.New().String()
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return nil, err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: s.table, Item: item,
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (s *DynamoDBTaskStore) Update(ctx context.Context, id string, params Update) error {
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
func (s *DynamoDBTaskStore) Delete(ctx context.Context, id string) error {
	key, err := GetKey(id)
	if err != nil {
		return err
	}
	res, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:    s.table,
		Key:          key,
		ReturnValues: ReturnAllOld,
	})
	if err != nil {
		return err
	}
	if len(res.Attributes) == 0 {
		return ErrorNotFound
	}
	return nil
}
func (s *DynamoDBTaskStore) GetTaskByID(ctx context.Context, id string) (*types.Task, error) {
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
	var task *types.Task
	if err := attributevalue.UnmarshalMap(res.Item, &task); err != nil {
		return nil, ErrorNotFound
	}
	return task, nil
}

// TODO: Pagination
// Filter
func (s *DynamoDBTaskStore) GetTasks(ctx context.Context, filter Filter, pagination *Pagination) ([]*types.Task, error) {
	var tasks []*types.Task
	numPages := 0
	limit := int32(pagination.Limit)
	page := int(pagination.Page)
	paginator := dynamodb.NewScanPaginator(s.client, &dynamodb.ScanInput{
		TableName: s.table,
		Limit:     &limit,
	})
	for paginator.HasMorePages() {
		res, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if numPages == page {
			if err := attributevalue.UnmarshalListOfMaps(res.Items, &tasks); err != nil {
				return nil, err
			}
			return tasks, nil
		}
		numPages++
	}
	return nil, ErrorNotFound
}
func (s *DynamoDBTaskStore) Drop(ctx context.Context) error {
	_, err := s.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: s.table,
	})
	return err
}
