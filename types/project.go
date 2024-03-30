package types

import "fmt"

type Project struct {
	ID          string   `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	UserID      string   `bson:"userID" json:"userID"`
	Tasks       []string `bson:"tasks" json:"tasks"`
}

const (
	minProjectTitleLen       = 5
	minProjectDescriptionLen = 5
)

type CreateProjectParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewProjectFromParams(params CreateProjectParams) *Project {
	return &Project{
		Title:       params.Title,
		Description: params.Description,
	}
}
func (params CreateProjectParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Title) < minProjectTitleLen {
		errors["title"] = fmt.Sprintf("Title length should be at least %d", minProjectTitleLen)
	}
	if len(params.Description) < minProjectTitleLen {
		errors["description"] = fmt.Sprintf("Description length should be at least %d", minProjectDescriptionLen)
	}
	return errors
}
