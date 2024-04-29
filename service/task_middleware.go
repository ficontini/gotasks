package service

import (
	"context"
	"time"

	"github.com/ficontini/gotasks/types"
	"github.com/sirupsen/logrus"
)

type TaskLogMiddleware struct {
	next TaskServicer
}

func NewTaskLogMiddleware(next TaskServicer) TaskServicer {
	return &TaskLogMiddleware{
		next: next,
	}
}
func (m *TaskLogMiddleware) GetTaskByID(ctx context.Context, id string) (task *types.Task, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   id,
			"err":  err,
		}).Info("Get task by ID")
	}(time.Now())
	task, err = m.next.GetTaskByID(ctx, id)
	return task, err
}
func (m *TaskLogMiddleware) CreateTask(ctx context.Context, params types.NewTaskParams) (task *types.Task, err error) {
	defer func(start time.Time) {
		var (
			id    string
			title string
		)
		if task != nil {
			id = task.ID
			title = task.Title
		}
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"id":    id,
			"title": title,
			"err":   err,
		}).Info("Create task")
	}(time.Now())
	task, err = m.next.CreateTask(ctx, params)
	return task, err

}
func (m *TaskLogMiddleware) GetTasks(ctx context.Context, params *TaskQueryParams) (tasks []*types.Task, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Get tasks")
	}(time.Now())
	tasks, err = m.next.GetTasks(ctx, params)
	return tasks, err

}
func (m *TaskLogMiddleware) GetTasksByUserID(ctx context.Context, id string, params TaskQueryParams) (tasks []*types.Task, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"userID": id,
			"err":    err,
		}).Info("Get tasks by user")
	}(time.Now())
	tasks, err = m.next.GetTasksByUserID(ctx, id, params)
	return tasks, err
}
func (m *TaskLogMiddleware) DeleteTask(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"taskID": id,
			"err":    err,
		}).Info("Delete task")
	}(time.Now())
	err = m.next.DeleteTask(ctx, id)
	return err

}
func (m *TaskLogMiddleware) CompleteTask(ctx context.Context, params types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"taskID": params.TaskID,
			"userID": params.UserID,
			"err":    err,
		}).Info("Complete task")
	}(time.Now())
	err = m.next.CompleteTask(ctx, params)
	return err
}
func (m *TaskLogMiddleware) AssignTaskToUser(ctx context.Context, req types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"taskID": req.TaskID,
			"userID": req.UserID,
			"err":    err,
		}).Info("Assign task to user")
	}(time.Now())
	err = m.next.AssignTaskToUser(ctx, req)
	return err
}
func (m *TaskLogMiddleware) AssignTaskToSelf(ctx context.Context, req types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"taskID": req.TaskID,
			"userID": req.UserID,
			"err":    err,
		}).Info("Assign task to user")
	}(time.Now())
	err = m.next.AssignTaskToSelf(ctx, req)
	return err
}

func (m *TaskLogMiddleware) UpdateDueDate(ctx context.Context, id string, params types.UpdateDueDateTaskRequest) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"taskID": id,
			"err":    err,
		}).Info("Update due date")
	}(time.Now())
	err = m.next.UpdateDueDate(ctx, id, params)
	return err
}
