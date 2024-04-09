package types

import "go.mongodb.org/mongo-driver/bson"

type Filter interface {
	Apply() bson.M
}
type EmptyFilter struct{}

func (f EmptyFilter) Apply() bson.M {
	return map[string]any{}
}

type CompletedFilter struct {
	Completed bool
}

func (f CompletedFilter) Apply() bson.M {
	return bson.M{"completed": f.Completed}
}

type AssignedToFilter struct {
	AssignedTo interface{}
}

func (f AssignedToFilter) Apply() bson.M {
	return bson.M{"assignedTo": f.AssignedTo}
}

type UserTasksFilter struct {
	CompletedFilter
	AssignedToFilter
}

func (f UserTasksFilter) Apply() bson.M {
	return bson.M{
		"$and": []bson.M{
			f.CompletedFilter.Apply(),
			f.AssignedToFilter.Apply(),
		},
	}
}
