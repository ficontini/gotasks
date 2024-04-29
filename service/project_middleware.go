package service

import (
	"context"
	"time"

	"github.com/ficontini/gotasks/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next ProjectServicer
}

func NewLogMiddleware(next ProjectServicer) ProjectServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CreateProject(ctx context.Context, params types.NewProjectParams, userID string) (project *types.Project, err error) {
	defer func(start time.Time) {
		var (
			projectID string
			title     string
		)
		if project != nil {
			projectID = project.ID
			title = project.Title
		}
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"projectID": projectID,
			"title":     title,
			"userID":    userID,
			"err":       err,
		}).Info("Create project")
	}(time.Now())
	project, err = m.next.CreateProject(ctx, params, userID)
	return project, err
}
func (m *LogMiddleware) GetProjectByID(ctx context.Context, id string) (project *types.Project, err error) {
	defer func(start time.Time) {
		var (
			title  string
			userID string
		)
		if project != nil {
			title = project.Title
			userID = project.UserID
		}
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"projectID": id,
			"title":     title,
			"userID":    userID,
			"err":       err,
		}).Info("Get project by id")
	}(time.Now())
	project, err = m.next.GetProjectByID(ctx, id)
	return project, err

}
func (m *LogMiddleware) AddTask(ctx context.Context, projectID string, params types.AddTaskParams) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"projectID": projectID,
			"taskID":    params.TaskID,
			"err":       err,
		}).Info("Add task")
	}(time.Now())
	err = m.next.AddTask(ctx, projectID, params)
	return err
}
