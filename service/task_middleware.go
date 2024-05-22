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
		if err != nil {
			logrus.WithError(err).Error("Failed to get task")
		} else {
			logrus.WithFields(logrus.Fields{
				"taskID": task.ID,
				"took":   time.Since(start),
			}).Info("Task retrieved successfully")
		}
	}(time.Now())
	task, err = m.next.GetTaskByID(ctx, id)
	return task, err
}
func (m *TaskLogMiddleware) CreateTask(ctx context.Context, params types.NewTaskParams) (task *types.Task, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to create task")
		} else {
			logrus.WithFields(logrus.Fields{
				"taskID": task.ID,
				"took":   time.Since(start),
			}).Info("Task created successfully")
		}
	}(time.Now())
	task, err = m.next.CreateTask(ctx, params)
	return task, err

}
func (m *TaskLogMiddleware) GetTasks(ctx context.Context, params *TaskQueryParams) (tasks []*types.Task, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get tasks")
		} else {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
			}).Info("Get tasks")
		}
	}(time.Now())
	tasks, err = m.next.GetTasks(ctx, params)
	return tasks, err

}
func (m *TaskLogMiddleware) GetTasksByUserID(ctx context.Context, id string, params TaskQueryParams) (tasks []*types.Task, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get tasks by user ID")
		} else {
			logrus.WithFields(logrus.Fields{
				"userID": id,
				"took":   time.Since(start),
			}).Info("Get tasks")
		}
	}(time.Now())
	tasks, err = m.next.GetTasksByUserID(ctx, id, params)
	return tasks, err
}
func (m *TaskLogMiddleware) DeleteTask(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to delete task")
		} else {
			logrus.WithFields(logrus.Fields{
				"taskID": id,
				"took":   time.Since(start),
			}).Info("DeleteTask successfully completed")
		}
	}(time.Now())
	err = m.next.DeleteTask(ctx, id)
	return err

}
func (m *TaskLogMiddleware) CompleteTask(ctx context.Context, params types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to complete task")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"taskID": params.TaskID,
			}).Info("CompleteTask successfully completed")
		}
	}(time.Now())
	err = m.next.CompleteTask(ctx, params)
	return err
}
func (m *TaskLogMiddleware) AssignTaskToUser(ctx context.Context, req types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to assign task to user")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": req.UserID,
				"taskID": req.TaskID,
			}).Info("AssignTaskToUser successfully completed")
		}
	}(time.Now())
	err = m.next.AssignTaskToUser(ctx, req)
	return err
}
func (m *TaskLogMiddleware) AssignTaskToSelf(ctx context.Context, req types.UpdateTaskRequest) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to assign task to user")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": req.UserID,
				"taskID": req.TaskID,
			}).Info("AssignTaskToSelf successfully completed")
		}
	}(time.Now())
	err = m.next.AssignTaskToSelf(ctx, req)
	return err
}

func (m *TaskLogMiddleware) UpdateDueDate(ctx context.Context, id string, params types.UpdateDueDateTaskRequest) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to update task due date")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"taskID": id,
			}).Info("UpdateDueDate successfully completed")
		}
	}(time.Now())
	err = m.next.UpdateDueDate(ctx, id, params)
	return err
}
