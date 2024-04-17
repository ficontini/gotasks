package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

type ProjectService struct {
	store *db.Store
}

func NewProjectService(store *db.Store) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (svc *ProjectService) CreateProject(ctx context.Context, params types.NewProjectParams, userID string) (*types.Project, error) {
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

// TODO:
func (svc *ProjectService) AddTask(ctx context.Context, params types.AddTaskParams) error {
	exists, task := svc.taskExists(ctx, params.TaskID)
	if !exists {
		return ErrTaskNotFound
	}
	if !svc.projectExists(ctx, params.ProjectID) {
		return ErrProjectNotFound
	}
	if task.ProjectID == params.ProjectID {
		return ErrTaskAlreadyAssociated
	}
	if err := svc.store.Project.UpdateProjectTasks(ctx, params); err != nil {
		return err
	}
	return nil
}
func (svc *ProjectService) taskExists(ctx context.Context, id string) (bool, *types.Task) {
	task, err := svc.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		return false, nil
	}
	return task != nil, task
}
func (svc *ProjectService) projectExists(ctx context.Context, id string) bool {
	project, err := svc.store.Project.GetProjectByID(ctx, id)
	if err != nil {
		return false
	}
	return project != nil
}

var (
	ErrTaskAlreadyAssociated = errors.New("task is already associated with this project")
	ErrProjectNotFound       = errors.New("project resource not found")
)
