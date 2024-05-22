package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
)

type Filter interface {
	ToBSON() bson.M
	ToExpression() (expression.Expression, error)
}
type EmptyFilter struct {
	DataType DataTyper
}

func NewEmptyFilter(dataType DataTyper) Filter {
	return &EmptyFilter{
		DataType: dataType,
	}
}
func (f EmptyFilter) ToBSON() bson.M {
	return bson.M{}
}
func (f EmptyFilter) ToExpression() (expression.Expression, error) {
	KeyCond := f.DataType.GetKeyCondition()
	return expression.NewBuilder().WithKeyCondition(KeyCond).Build()
}

type SimpleFilter struct {
	DataType DataTyper
	Field    FieldFilterer
}

func NewSimpleFilter(dataType DataTyper, field FieldFilterer) Filter {
	return &SimpleFilter{
		DataType: dataType,
		Field:    field,
	}
}
func (f SimpleFilter) ToBSON() bson.M {
	return f.Field.GetBSONFilter()
}

func (f SimpleFilter) ToExpression() (expression.Expression, error) {
	return buildExpression(f.DataType, f.Field.GetFilter())
}

type CompositeFilter struct {
	DataType DataTyper
	Field1   FieldFilterer
	Field2   FieldFilterer
}

func NewCompositeFilter(dataType DataTyper, field1, field2 FieldFilterer) Filter {
	return &CompositeFilter{
		DataType: dataType,
		Field1:   field1,
		Field2:   field2,
	}
}
func (f CompositeFilter) ToBSON() bson.M {
	return bson.M{"$and": []bson.M{
		f.Field1.GetBSONFilter(),
		f.Field2.GetBSONFilter(),
	}}
}

func (f CompositeFilter) ToExpression() (expression.Expression, error) {
	filter := expression.And(f.Field1.GetFilter(), f.Field2.GetFilter())
	return buildExpression(f.DataType, filter)

}

func buildExpression(dataType DataTyper, filter expression.ConditionBuilder) (expression.Expression, error) {
	keyCond := dataType.GetKeyCondition()
	return expression.NewBuilder().WithFilter(filter).WithKeyCondition(keyCond).Build()
}
