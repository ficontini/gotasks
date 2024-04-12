package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

type TaskService struct {
	store *db.Store
}

func NewTaskService(store *db.Store) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (svc *TaskService) GetTaskByID(ctx context.Context, id string) (*types.Task, error) {
	task, err := svc.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) CreateTask(ctx context.Context, params types.NewTaskParams) (*types.Task, error) {
	task := types.NewTaskFromParams(params)
	insertedTask, err := svc.store.Task.InsertTask(ctx, task)
	return insertedTask, err
}

type TaskQueryParams struct {
	db.Pagination
	Completed *bool
}

func (svc *TaskService) GetTasks(ctx context.Context, params *TaskQueryParams) ([]*types.Task, error) {
	return svc.store.Task.GetTasks(ctx, db.NewCompletedFilter(params.Completed), &params.Pagination)
}

func (svc *TaskService) GetTasksByUserID(ctx context.Context, id string, params TaskQueryParams) ([]*types.Task, error) {
	filter, err := db.NewUserTasksFilter(params.Completed, id)
	if err != nil {
		return nil, err
	}
	return svc.store.Task.GetTasks(ctx, filter, &params.Pagination)
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
func (svc *TaskService) CompleteTask(ctx context.Context, params types.UpdateTaskRequest) error {
	task, err := svc.getTask(ctx, params.TaskID)
	if err != nil {
		return err
	}

	if task.AssignedTo != params.UserID {
		return ErrUnAuthorized
	}
	if task.Completed {
		return ErrTaskAlreadyCompleted
	}
	update := db.TaskCompleteUpdater{Completed: true}
	return svc.store.Task.Update(ctx, task.ID, update)
}
func (svc *TaskService) getTask(ctx context.Context, id string) (*types.Task, error) {
	task, err := svc.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}
func (svc *TaskService) AssignTaskToSelf(ctx context.Context, req types.UpdateTaskRequest) error {
	return svc.assignTask(ctx, req.TaskID, req.UserID)
}
func (svc *TaskService) AssignTaskToUser(ctx context.Context, req types.UpdateTaskRequest) error {
	if _, err := svc.store.User.GetUserByID(ctx, req.UserID); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrUserNotFound
		}
	}
	return svc.assignTask(ctx, req.TaskID, req.UserID)
}
func (svc *TaskService) assignTask(ctx context.Context, taskID, userID string) error {
	if err := svc.store.Task.Update(ctx, taskID, db.TaskAssignationUpdater{AssignedTo: userID}); err != nil {
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
