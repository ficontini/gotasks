package db

import "go.mongodb.org/mongo-driver/bson"

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
	return bson.M{"completed": f.Completed}
}

type AssignedToFilter struct {
	AssignedTo interface{}
}

func (f AssignedToFilter) ToBSON() bson.M {
	return bson.M{"assignedTo": f.AssignedTo}
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
