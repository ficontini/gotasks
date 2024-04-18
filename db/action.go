package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBAction interface {
	get() (interface{}, error)
}
type UpdateAction struct {
	Key       map[string]dynamodbtypes.AttributeValue
	Update    expression.UpdateBuilder
	TableName string
}

func NewTaskUpdateAction(id string, params Update) (*UpdateAction, error) {
	key, err := GetKey(id)
	if err != nil {
		return nil, err
	}
	return &UpdateAction{
		Key:       key,
		Update:    params.ToExpression(),
		TableName: taskColl,
	}, nil
}
func NewProjectUpdateAction(id string, params Update) (*UpdateAction, error) {
	key, err := GetKey(id)
	if err != nil {
		return nil, err
	}
	return &UpdateAction{
		Key:       key,
		Update:    params.ToExpression(),
		TableName: projectColl,
	}, nil
}
func (a *UpdateAction) get() (interface{}, error) {
	expr, err := expression.NewBuilder().WithUpdate(a.Update).Build()
	if err != nil {
		return nil, err
	}
	return &dynamodbtypes.Update{
		TableName:                 &a.TableName,
		Key:                       a.Key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}, nil
}
