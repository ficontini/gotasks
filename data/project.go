package data

import "fmt"

type Project struct {
	ID          ID     `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	UserID      ID     `bson:"userID" json:"userID"`
	Tasks       []ID   `bson:"tasks" json:"tasks"`
}

func (project *Project) ContainsTask(taskID ID) bool {
	for _, id := range project.Tasks {
		if id == taskID {
			return true
		}
	}
	return false
}

type CreateProjectParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewProjectFromParams(params CreateProjectParams) *Project {
	return &Project{
		Title:       params.Title,
		Description: params.Description,
		Tasks:       []ID{},
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

type AddTaskParams struct {
	TaskID string `json:"taskID"`
}

func (p AddTaskParams) ToMap() map[string]any {
	return map[string]any{
		"tasks": ID(p.TaskID),
	}
}
