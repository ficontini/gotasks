package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/gotasks/types"
	"github.com/google/uuid"
)

const (
	ReturnAllOld = "ALL_OLD"
)

type DynamoDBTaskStore struct {
	client *dynamodb.Client
	table  *string
	gsi    *string
}

func NewDynamoDBTaskStore(client *dynamodb.Client) *DynamoDBTaskStore {
	return &DynamoDBTaskStore{
		client: client,
		table:  aws.String(taskColl),
		gsi:    aws.String(dataTypeGSI),
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

// TODO: review
func (s *DynamoDBTaskStore) GetTasks(ctx context.Context, filter Filter, pagination *Pagination) ([]*types.Task, error) {
	expr, err := filter.ToExpression()
	if err != nil {
		return nil, err
	}
	pagination.generatePaginationForDynamoDB()

	queryInput := &dynamodb.QueryInput{
		TableName:                 s.table,
		IndexName:                 s.gsi,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		Limit:                     aws.Int32(int32(pagination.Limit)),
	}
	opts := NewDynamoDBQueryOptions(queryInput, pagination)
	startT := time.Now()
	collectiveResult, err := PaginatedDynamoDBQuery(ctx, s.client, opts)
	fmt.Println("TIME::::::", time.Since(startT).Seconds())
	if err != nil {
		return nil, err
	}
	start := pagination.Offset
	var tasks []*types.Task
	if start > len(collectiveResult) {
		return tasks, nil
	}
	endIdx := min(start+int(pagination.Limit), len(collectiveResult))
	if err := attributevalue.UnmarshalListOfMaps(collectiveResult[start:endIdx], &tasks); err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *DynamoDBTaskStore) Drop(ctx context.Context) error {
	_, err := s.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: s.table,
	})
	return err
}
