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

func (svc *ProjectService) CreateProject(ctx context.Context, params data.CreateProjectParams, userID data.ID) (*data.Project, error) {
	project := data.NewProjectFromParams(params)
	project.UserID = userID
	return svc.store.Project.InsertProject(ctx, project)
}

func (svc *ProjectService) GetProjectByID(ctx context.Context, id string) (*data.Project, error) {
	return svc.store.Project.GetProjectByID(ctx, data.ID(id))
}
func (svc *ProjectService) getTaskByID(ctx context.Context, id string) (*data.Task, error) {
	return svc.store.Task.GetTaskByID(ctx, data.ID(id))
}
func (svc *ProjectService) GetTasksByProject(ctx context.Context, id string, pagination *db.Pagination) ([]*data.Task, error) {
	return svc.store.Task.GetTasksByProject(ctx, data.ID(id), pagination)
}
func (svc *ProjectService) AddTask(ctx context.Context, project *data.Project, taskID string) error {
	task, err := svc.getTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	if project.ContainsTask(task.ID) {
		return ErrTaskAlreadyAssociated
	}
	return svc.store.Project.UpdateProjectTasks(ctx, project.ID, task.ID)
}

var (
	ErrTaskAlreadyAssociated = errors.New("task is already associated with this project")
)
