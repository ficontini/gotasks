package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
)

type TaskService struct {
	taskStore db.TaskStore
}

func NewTaskService(taskStore db.TaskStore) *TaskService {
	return &TaskService{
		taskStore: taskStore,
	}
}

func (svc *TaskService) GetTaskByID(ctx context.Context, id string) (*data.Task, error) {
	task, err := svc.taskStore.GetTaskByID(ctx, data.ID(id))
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) CreateTask(ctx context.Context, params data.CreateTaskParams) (*data.Task, error) {
	task := data.NewTaskFromParams(params)
	insertedTask, err := svc.taskStore.InsertTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return insertedTask, nil
}

// TODO: review
func createCompletionFilter(completed *bool) db.Map {
	if completed != nil {
		return db.Map{"completed": completed}
	} else {
		return nil
	}
}

func (svc *TaskService) GetTasks(ctx context.Context, params *db.TaskQueryParams) ([]*data.Task, error) {
	return svc.taskStore.GetTasks(ctx, createCompletionFilter(params.Completed), &params.Pagination)
}

func (svc *TaskService) DeleteTask(ctx context.Context, id string) error {
	return svc.taskStore.Delete(ctx, data.ID(id))
}
func (svc *TaskService) CompleteTask(ctx context.Context, idStr string) error {
	id := data.ID(idStr)
	completed, err := svc.IsTaskCompleted(ctx, id)
	if err != nil {
		return err
	}
	if completed {
		return ErrTaskAlreadyCompleted
	}
	params := data.UpdateTaskParams{
		Completed: true,
	}
	return svc.taskStore.Update(ctx, data.ID(id), params)
}
func (svc *TaskService) IsTaskCompleted(ctx context.Context, id data.ID) (bool, error) {
	task, err := svc.taskStore.GetTaskByID(ctx, id)
	if err != nil {
		return false, err
	}
	return task.Completed, nil
}

var (
	ErrTaskAlreadyCompleted = errors.New("task already completed")
)
