package db

func NewCompletedFilter(completed *bool) Filter {
	if completed != nil {
		return CompletedFilter{Completed: *completed}
	} else {
		return EmptyFilter{}
	}
}
func NewUserTasksFilter(completed *bool, id string) (Filter, error) {
	assignedToFilter := AssignedToFilter{AssignedTo: id}
	if completed != nil {
		return UserTasksFilter{
			CompletedFilter:  CompletedFilter{Completed: *completed},
			AssignedToFilter: assignedToFilter,
		}, nil

	}
	return assignedToFilter, nil
}
