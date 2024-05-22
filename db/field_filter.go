package db

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataTyper interface {
	GetKeyCondition() expression.KeyConditionBuilder
}
type DataType struct {
	DataType string
}

func NewDataType(dataType string) DataTyper {
	return &DataType{
		DataType: dataType,
	}
}
func (d *DataType) GetKeyCondition() expression.KeyConditionBuilder {
	return expression.Key(dataTypeField).Equal(expression.Value(d.DataType))
}

type FieldFilterer interface {
	GetBSONFilter() bson.M
	GetFilter() expression.ConditionBuilder
}
type CompleteFieldFilterer struct {
	Completed bool
}

func NewCompletedFieldFilterer(completed bool) FieldFilterer {
	return &CompleteFieldFilterer{
		Completed: completed,
	}
}
func (c *CompleteFieldFilterer) GetBSONFilter() bson.M {
	return bson.M{completedField: c.Completed}
}
func (c *CompleteFieldFilterer) GetFilter() expression.ConditionBuilder {
	return expression.Equal(expression.Name(completedField), expression.Value(c.Completed))
}

type AssigneFieldFilterer struct {
	AssignedTo string
}

func NewAssigneFieldFilterer(assignedTo string) FieldFilterer {
	return &AssigneFieldFilterer{
		AssignedTo: assignedTo,
	}
}
func (c *AssigneFieldFilterer) GetBSONFilter() bson.M {
	//TODO: Review how to handle this error
	oid, err := primitive.ObjectIDFromHex(c.AssignedTo)
	if err != nil {
		log.Fatal(err)
	}
	return bson.M{assignedToField: oid}
}
func (c *AssigneFieldFilterer) GetFilter() expression.ConditionBuilder {
	return expression.Equal(expression.Name(assignedToField), expression.Value(c.AssignedTo))
}
