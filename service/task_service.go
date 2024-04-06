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
	task, err := svc.taskStore.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) CreateTask(ctx context.Context, params data.CreateTaskParams) (*data.Task, error) {
	task := data.NewTaskFromParams(params)
	insertedTask, err := svc.taskStore.InsertTask(ctx, task)
	return insertedTask, err
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
	if err := svc.taskStore.Delete(ctx, id); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound
		}
		return err
	}
	return nil
}
func (svc *TaskService) CompleteTask(ctx context.Context, id string) error {
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
	return svc.taskStore.Update(ctx, id, params)
}
func (svc *TaskService) IsTaskCompleted(ctx context.Context, id string) (bool, error) {
	task, err := svc.taskStore.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return false, ErrResourceNotFound
		}
		return false, err
	}
	return task.Completed, nil
}

var (
	ErrTaskAlreadyCompleted = errors.New("task already completed")
)
