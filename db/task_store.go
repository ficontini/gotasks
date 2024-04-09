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
	Deleter
}

type TaskUpdater interface {
	//Update(context.Context, string, Map) error
	SetTaskAsComplete(context.Context, string, types.TaskCompletionRequest) error
	SetTaskAssignee(context.Context, string, types.TaskAssignmentRequest) error
}
type TaskGetter interface {
	GetTasks(context.Context, types.Filter, *Pagination) ([]*types.Task, error)
	GetTaskByID(context.Context, string) (*types.Task, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}
