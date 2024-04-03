package db

import (
	"context"

	"github.com/ficontini/gotasks/data"
)

const taskColl = "tasks"

type TaskStore interface {
	InsertTask(context.Context, *data.Task) (*data.Task, error)
	TaskUpdater
	TaskGetter
	Deleter
}

type TaskUpdater interface {
	Update(context.Context, data.ID, data.UpdateTaskParams) error
	UpdateTaskProjects(context.Context, data.ID, data.ID) error
}
type TaskGetter interface {
	GetTasks(context.Context, Map, *Pagination) ([]*data.Task, error)
	GetTasksByProject(context.Context, data.ID, *Pagination) ([]*data.Task, error)
	GetTaskByID(context.Context, data.ID) (*data.Task, error)
}
