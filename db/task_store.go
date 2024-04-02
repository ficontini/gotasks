package db

import (
	"context"

	"github.com/ficontini/gotasks/types"
)

const taskColl = "tasks"

type TaskStore interface {
	GetTaskByID(context.Context, types.ID) (*types.Task, error)
	GetTasks(context.Context, Map, *Pagination) ([]*types.Task, error)
	InsertTask(context.Context, *types.Task) (*types.Task, error)
	Updater
	Delete(context.Context, types.ID) error
}
