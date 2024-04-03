package db

import (
	"context"

	"github.com/ficontini/gotasks/types"
)

const taskColl = "tasks"

type TaskStore interface {
	InsertTask(context.Context, *types.Task) (*types.Task, error)
	TaskUpdater
	TaskGetter
	Delete(context.Context, types.ID) error
}

type TaskUpdater interface {
	Update(context.Context, types.ID, types.UpdateTaskParams) error
	UpdateTaskProjects(context.Context, types.ID, types.ID) error
}
type TaskGetter interface {
	GetTasks(context.Context, Map, *Pagination) ([]*types.Task, error)
	GetTasksByProject(context.Context, types.ID, *Pagination) ([]*types.Task, error)
	GetTaskByID(context.Context, types.ID) (*types.Task, error)
}
