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
	//Update(context.Context, string, Map) error
	SetTaskAsComplete(context.Context, string, data.TaskCompletionRequest) error
	SetTaskAssignee(context.Context, string, data.TaskAssignmentRequest) error
}
type TaskGetter interface {
	GetTasks(context.Context, data.Filter, *Pagination) ([]*data.Task, error)
	GetTaskByID(context.Context, string) (*data.Task, error)
	GetTasksByUserID(context.Context, data.Filter, *Pagination) ([]*data.Task, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}
