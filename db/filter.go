package db

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter interface {
	ToBSON() bson.M
	ToExpression() (expression.Expression, error)
}
type EmptyFilter struct{}

func (f EmptyFilter) ToBSON() bson.M {
	return map[string]any{}
}

// TODO: Review
func (f EmptyFilter) ToExpression() (expression.Expression, error) {
	return expression.Expression{}, nil
}

type CompletedFilter struct {
	Completed bool
}

func (f CompletedFilter) ToBSON() bson.M {
	return bson.M{completedField: f.Completed}
}
func (f CompletedFilter) getConditionBuilder() expression.ConditionBuilder {
	return expression.Equal(expression.Name(completedField), expression.Value(f.Completed))
}
func (f CompletedFilter) ToExpression() (expression.Expression, error) {
	filter := expression.Equal(expression.Name(completedField), expression.Value(f.Completed))
	return expression.NewBuilder().WithFilter(filter).Build()
}

type AssignedToFilter struct {
	AssignedTo string
}

// TODO: Review
func (f AssignedToFilter) ToBSON() bson.M {
	oid, err := primitive.ObjectIDFromHex(f.AssignedTo)
	if err != nil {
		log.Fatal(err)
	}
	return bson.M{assignedToField: oid}
}
func (f AssignedToFilter) getConditionBuilder() expression.ConditionBuilder {
	return expression.Equal(expression.Name(assignedToField), expression.Value(f.AssignedTo))
}
func (f AssignedToFilter) ToExpression() (expression.Expression, error) {
	KeyCond := expression.Key(assignedToField).Equal(expression.Value(f.AssignedTo))
	return expression.NewBuilder().WithKeyCondition(KeyCond).Build()
}

type UserTasksFilter struct {
	CompletedFilter
	AssignedToFilter
}

func (f UserTasksFilter) ToBSON() bson.M {
	return bson.M{
		"$and": []bson.M{
			f.CompletedFilter.ToBSON(),
			f.AssignedToFilter.ToBSON(),
		},
	}
}

// TODO:
func (f UserTasksFilter) ToExpression() (expression.Expression, error) {
	filter := expression.And(f.CompletedFilter.getConditionBuilder(), f.AssignedToFilter.getConditionBuilder())
	return expression.NewBuilder().WithFilter(filter).Build()
}
