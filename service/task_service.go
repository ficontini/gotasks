package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
)

type TaskService struct {
	store *db.Store
}

func NewTaskService(store *db.Store) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (svc *TaskService) GetTaskByID(ctx context.Context, id string) (*data.Task, error) {
	task, err := svc.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) CreateTask(ctx context.Context, params data.CreateTaskParams) (*data.Task, error) {
	task := data.NewTaskFromParams(params)
	insertedTask, err := svc.store.Task.InsertTask(ctx, task)
	return insertedTask, err
}

type TaskQueryParams struct {
	db.Pagination
	Completed *bool
}

// TODO: review
func createCompletionFilter(completed *bool) data.Filter {
	if completed != nil {
		return data.CompletionFilter{Completed: *completed}
	} else {
		return data.NoFilter{}
	}
}
func (svc *TaskService) GetTasks(ctx context.Context, params *TaskQueryParams) ([]*data.Task, error) {
	return svc.store.Task.GetTasks(ctx, createCompletionFilter(params.Completed), &params.Pagination)
}

func (svc *TaskService) GetTasksByUserID(ctx context.Context, id string, params TaskQueryParams) ([]*data.Task, error) {
	filter := data.AssignationFilter{AssignedTo: id}
	return svc.store.Task.GetTasksByUserID(ctx, filter, &params.Pagination)
}
func (svc *TaskService) DeleteTask(ctx context.Context, id string) error {
	if err := svc.store.Task.Delete(ctx, id); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}
func (svc *TaskService) CompleteTask(ctx context.Context, id, userID string) error {
	task, err := svc.getTask(ctx, id)
	if err != nil {
		return err
	}

	if task.AssignedTo != userID {
		return ErrUnAuthorized
	}
	if task.Completed {
		return ErrTaskAlreadyCompleted
	}
	params := data.TaskCompletionRequest{
		Completed: true,
	}
	return svc.store.Task.SetTaskAsComplete(ctx, id, params)
}
func (svc *TaskService) getTask(ctx context.Context, id string) (*data.Task, error) {
	task, err := svc.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) AssignMeTask(ctx context.Context, id string, req data.TaskAssignmentRequest) error {
	return svc.assignTask(ctx, id, req)
}
func (svc *TaskService) AssignTaskToUser(ctx context.Context, id string, req data.TaskAssignmentRequest) error {
	if _, err := svc.store.User.GetUserByID(ctx, req.UserID); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrUserNotFound
		}
	}
	return svc.assignTask(ctx, id, req)
}
func (svc *TaskService) assignTask(ctx context.Context, id string, req data.TaskAssignmentRequest) error {
	if err := svc.store.Task.SetTaskAssignee(ctx, id, req); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

var (
	ErrTaskAlreadyCompleted = errors.New("task already completed")
	ErrTaskNotFound         = errors.New("task resource not found")
	ErrUnAuthorized         = errors.New("unauthorized request")
)
