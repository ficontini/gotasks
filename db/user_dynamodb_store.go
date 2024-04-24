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

const emailGSI = "EmailGSI"
const dataTypeGSI = "dataTypeGSI"

type DynamoDBUserStore struct {
	client   *dynamodb.Client
	table    *string
	emailGSI *string
	queryGSI *string
}

func NewDynamoDBUserStore(client *dynamodb.Client) *DynamoDBUserStore {
	return &DynamoDBUserStore{
		client:   client,
		table:    aws.String(userColl),
		emailGSI: aws.String(emailGSI),
		queryGSI: aws.String(dataTypeGSI),
	}
}
func (s *DynamoDBUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	user.ID = uuid.New().String()
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: s.table, Item: item,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *DynamoDBUserStore) GetUserByID(ctx context.Context, idStr string) (*types.User, error) {
	key, err := GetKey(idStr)
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
	var user *types.User
	if err := attributevalue.UnmarshalMap(res.Item, &user); err != nil {
		return nil, ErrorNotFound
	}
	return user, nil
}
func (s *DynamoDBUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	keyEx := expression.Key(emailField).Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}
	queryOutput, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 s.table,
		IndexName:                 s.emailGSI,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return nil, ErrorNotFound
	}
	var user *types.User
	if err := attributevalue.UnmarshalMap(queryOutput.Items[0], &user); err != nil {
		return nil, ErrorNotFound
	}
	return user, nil
}

// TODO: Review
func (s *DynamoDBUserStore) GetUsers(ctx context.Context, filter Filter, pagination Pagination) ([]*types.User, error) {
	expr, err := filter.ToExpression()
	if err != nil {
		return nil, err
	}
	start, limit := pagination.generatePaginationForDynamoDB()
	queryInput := &dynamodb.QueryInput{
		TableName:                 s.table,
		IndexName:                 s.queryGSI,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(int32(limit)),
	}
	var collectiveResult []map[string]dynamodbtypes.AttributeValue
	paginator := dynamodb.NewQueryPaginator(s.client, queryInput)
	for {
		if !paginator.HasMorePages() {
			break
		}
		singlePage, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		collectiveResult = append(collectiveResult, singlePage.Items...)
		if len(collectiveResult) >= (start * limit) {
			break
		}
	}
	var users []*types.User
	if start > len(collectiveResult) {
		return users, nil
	}
	endIdx := min(start+limit, len(collectiveResult))
	if err := attributevalue.UnmarshalListOfMaps(collectiveResult[start:endIdx], &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *DynamoDBUserStore) Update(ctx context.Context, idStr string, params Update) error {
	key, err := GetKey(idStr)
	if err != nil {
		return err
	}
	update := params.ToExpression()
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	res, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
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
	//TODO: Review
	if len(res.Attributes) == 0 {
		return ErrorNotFound
	}
	return nil
}

func (s *DynamoDBUserStore) Drop(ctx context.Context) error {
	_, err := s.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: s.table,
	})
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
