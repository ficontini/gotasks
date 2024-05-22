package types

import (
	"fmt"
	"time"
)

const TaskDataType = "task"

type Task struct {
	ID          string    `bson:"_id,omitempty" dynamodbav:"ID" json:"id,omitempty"`
	Name        string    `bson:"name" dynamodbav:"name" json:"name"`
	Description string    `bson:"description,omitempty" dynamodbav:"description" json:"description,omitempty"`
	DueDate     time.Time `bson:"dueDate" dynamodbav:"dueDate" json:"dueDate"`
	Completed   bool      `bson:"completed" dynamodbav:"completed" json:"completed"`
	AssignedTo  string    `bson:"assignedTo" dynamodbav:"assignedTo" json:"assignedTo,omitempty"`
	ProjectID   string    `bson:"projectID" dynamodbav:"projectID" json:"projectID,omitempty"`
	DataType    string    `bson:"-" dynamodbav:"dataType" json:"-"`
}

type NewTaskParams struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

func NewTaskFromParams(params NewTaskParams) *Task {
	return &Task{
		Name:        params.Name,
		Description: params.Description,
		DueDate:     params.DueDate,
		DataType:    TaskDataType,
	}
}
func (params NewTaskParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Name) < minNameLen {
		errors["title"] = fmt.Sprintf("Title length should be at least %d", minNameLen)
	}
	if len(params.Description) < minDescriptionLen {
		errors["description"] = fmt.Sprintf("Description length should be at least %d", minDescriptionLen)
	}
	if !isDateValid(params.DueDate) {
		errors["dueDate"] = fmt.Sprintf("date %v is not valid", params.DueDate)
	}
	return errors
}
func isDateValid(date time.Time) bool {
	return date.After(time.Now())
}

type UpdateTaskRequest struct {
	TaskID string
	UserID string `json:"userID"`
}

type UpdateDueDateTaskRequest struct {
	DueDate    time.Time `json:"dueDate"`
	AssignedTo string
}

func (req UpdateDueDateTaskRequest) Validate() error {
	if !isDateValid(req.DueDate) {
		return fmt.Errorf("date %v is not valid", req.DueDate)
	}
	return nil
}
