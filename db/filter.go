package db

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter interface {
	ToBSON() bson.M
	ToExpression() expression.ConditionBuilder
}
type EmptyFilter struct{}

func (f EmptyFilter) ToBSON() bson.M {
	return map[string]any{}
}

// TODO: Review
func (f EmptyFilter) ToExpression() expression.ConditionBuilder {
	return expression.ConditionBuilder{}
}

type CompletedFilter struct {
	Completed bool
}

func (f CompletedFilter) ToBSON() bson.M {
	return bson.M{completedField: f.Completed}
}
func (f CompletedFilter) ToExpression() expression.ConditionBuilder {
	return expression.Equal(expression.Name(completedField), expression.Value(f.Completed))
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

// TODO: review
func (f AssignedToFilter) ToExpression() expression.ConditionBuilder {
	return expression.Equal(expression.Name(assignedToField), expression.Value(f.AssignedTo))
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
func (f UserTasksFilter) ToExpression() expression.ConditionBuilder {
	return expression.ConditionBuilder{}
}
