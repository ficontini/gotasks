package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBAction interface {
	get() (interface{}, error)
}
type UpdateAction struct {
	ID        string
	Params    Update
	TableName string
}

func NewTaskUpdateAction(id string, params Update) (*UpdateAction, error) {
	return &UpdateAction{
		ID:        id,
		Params:    params,
		TableName: taskColl,
	}, nil
}
func NewProjectUpdateAction(id string, params Update) (*UpdateAction, error) {
	return &UpdateAction{
		ID:        id,
		Params:    params,
		TableName: projectColl,
	}, nil
}
func (a *UpdateAction) get() (interface{}, error) {
	key, err := GetKey(a.ID)
	if err != nil {
		return nil, err
	}
	update := a.Params.ToExpression()
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, err
	}
	return &dynamodbtypes.Update{
		TableName:                 &a.TableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}, nil
}
