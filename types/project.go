package types

import "fmt"

type Project struct {
	ID          string   `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	UserID      string   `bson:"userID" json:"userID"`
	Tasks       []string `bson:"tasks" json:"tasks"`
}

type CreateProjectParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewProjectFromParams(params CreateProjectParams) *Project {
	return &Project{
		Title:       params.Title,
		Description: params.Description,
		Tasks:       []string{},
	}
}
func (params CreateProjectParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Title) < minTitleLen {
		errors["title"] = fmt.Sprintf("Title length should be at least %d", minTitleLen)
	}
	if len(params.Description) < minDescriptionLen {
		errors["description"] = fmt.Sprintf("Description length should be at least %d", minDescriptionLen)
	}
	return errors
}

type AddTaskRequest struct {
	TaskID string `json:"taskID"`
}
