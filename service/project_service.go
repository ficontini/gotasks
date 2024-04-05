package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
)

type ProjectService struct {
	store db.Store
}

func NewProjectService(store db.Store) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (svc *ProjectService) CreateProject(ctx context.Context, params data.CreateProjectParams, userID string) (*data.Project, error) {
	project := data.NewProjectFromParams(params)
	project.UserID = userID
	return svc.store.Project.InsertProject(ctx, project)
}

func (svc *ProjectService) GetProjectByID(ctx context.Context, id string) (*data.Project, error) {
	return svc.store.Project.GetProjectByID(ctx, id)
}

var (
	ErrTaskAlreadyAssociated = errors.New("task is already associated with this project")
)
