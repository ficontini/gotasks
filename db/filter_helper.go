package db

import "github.com/ficontini/gotasks/types"

func NewTaskCompletedFilter(completed *bool) Filter {
	dataType := NewDataType(types.TaskDataType)
	if completed != nil {
		return NewSimpleFilter(dataType, NewCompletedFieldFilterer(*completed))
	}
	return NewEmptyFilter(dataType)
}

func NewUserTasksFilter(completed *bool, id string) Filter {
	dataType := NewDataType(types.TaskDataType)
	fieldFiltered1 := NewAssigneFieldFilterer(id)
	if completed != nil {
		return NewCompositeFilter(dataType, fieldFiltered1, NewCompletedFieldFilterer(*completed))
	}
	return NewSimpleFilter(dataType, fieldFiltered1)
}
