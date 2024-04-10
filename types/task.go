package types

import (
	"fmt"
	"time"
)

type Task struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     time.Time `bson:"dueDate" json:"dueDate"`
	Completed   bool      `bson:"completed" json:"completed"`
	AssignedTo  string    `bson:"assignedTo" json:"assignedTo,omitempty"`
}

type NewTaskParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

func NewTaskFromParams(params NewTaskParams) *Task {
	return &Task{
		Title:       params.Title,
		Description: params.Description,
		DueDate:     params.DueDate,
	}
}
func (params NewTaskParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Title) < minDescriptionLen {
		errors["title"] = fmt.Sprintf("Title length should be at least %d", minTitleLen)
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
