package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter interface {
	ToBSON() bson.M
}
type EmptyFilter struct{}

func (f EmptyFilter) ToBSON() bson.M {
	return map[string]any{}
}

type CompletedFilter struct {
	Completed bool
}

func (f CompletedFilter) ToBSON() bson.M {
	return bson.M{completedField: f.Completed}
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
