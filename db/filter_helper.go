package db

import "go.mongodb.org/mongo-driver/bson/primitive"

func newAssignedToFilter(assignedTo string) (Filter, error) {
	oid, err := primitive.ObjectIDFromHex(assignedTo)
	if err != nil {
		return nil, err
	}
	return AssignedToFilter{AssignedTo: oid}, nil
}
func NewCompletedFilter(completed *bool) Filter {
	if completed != nil {
		return CompletedFilter{Completed: *completed}
	} else {
		return EmptyFilter{}
	}
}
func NewUserTasksFilter(completed *bool, id string) (Filter, error) {
	assignedToFilter, err := newAssignedToFilter(id)
	if err != nil {
		return nil, err
	}
	if completed != nil {
		return UserTasksFilter{
			CompletedFilter:  CompletedFilter{Completed: *completed},
			AssignedToFilter: assignedToFilter.(AssignedToFilter),
		}, nil

	}
	return assignedToFilter, nil
}
