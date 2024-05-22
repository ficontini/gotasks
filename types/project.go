package types

import "fmt"

type Project struct {
	ID          string   `bson:"_id,omitempty" dynamodbav:"ID" json:"id,omitempty"`
	Name        string   `bson:"name" dynamodbav:"name" json:"name"`
	Description string   `bson:"description" dynamodbav:"description" json:"description"`
	UserID      string   `bson:"userID" dynamodbav:"userID" json:"userID"`
	Tasks       []string `bson:"tasks" dynamodbav:"tasks" json:"tasks"`
}

func (project *Project) ContainsTask(taskID string) bool {
	for _, id := range project.Tasks {
		if id == taskID {
			return true
		}
	}
	return false
}

type NewProjectParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewProjectFromParams(params NewProjectParams) *Project {
	return &Project{
		Name:        params.Name,
		Description: params.Description,
		Tasks:       []string{},
	}
}
func (params NewProjectParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Name) < minNameLen {
		errors["title"] = fmt.Sprintf("Title length should be at least %d", minNameLen)
	}
	if len(params.Description) < minDescriptionLen {
		errors["description"] = fmt.Sprintf("Description length should be at least %d", minDescriptionLen)
	}
	return errors
}

type AddTaskParams struct {
	TaskID string `json:"taskID"`
}
