package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

type ProjectService struct {
	store db.Store
}

func NewProjectService(store db.Store) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (svc *ProjectService) CreateProject(ctx context.Context, params types.CreateProjectParams, userID string) (*types.Project, error) {
	project := types.NewProjectFromParams(params)
	project.UserID = userID
	return svc.store.Project.InsertProject(ctx, project)
}

func (svc *ProjectService) GetProjectByID(ctx context.Context, id string) (*types.Project, error) {
	project, err := svc.store.Project.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

var (
	ErrTaskAlreadyAssociated = errors.New("task is already associated with this project")
	ErrProjectNotFound       = errors.New("project resource not found")
)
