package db

import (
	"context"

	"github.com/ficontini/gotasks/types"
)

const taskColl = "tasks"

type TaskStore interface {
	TaskInserter
	TaskUpdater
	TaskGetter
	Deleter
	Dropper
}
type TaskInserter interface {
	InsertTask(context.Context, *types.Task) (*types.Task, error)
}
type TaskUpdater interface {
	Update(context.Context, string, Update) error
}
type TaskGetter interface {
	GetTasks(context.Context, Filter, *Pagination) ([]*types.Task, error)
	GetTaskByID(context.Context, string) (*types.Task, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}

type Dropper interface {
	Drop(context.Context) error
}
