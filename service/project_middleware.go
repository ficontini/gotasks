package service

import (
	"context"
	"time"

	"github.com/ficontini/gotasks/types"
	"github.com/sirupsen/logrus"
)

type ProjectLogMiddleware struct {
	next ProjectServicer
}

func NewProjectLogMiddleware(next ProjectServicer) ProjectServicer {
	return &ProjectLogMiddleware{
		next: next,
	}
}

func (m *ProjectLogMiddleware) CreateProject(ctx context.Context, params types.NewProjectParams, userID string) (project *types.Project, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to create project")
		} else {
			logrus.WithFields(logrus.Fields{
				"projectID": project.ID,
				"took":      time.Since(start),
			}).Info("Project created succesfully")
		}
	}(time.Now())
	project, err = m.next.CreateProject(ctx, params, userID)
	return project, err
}
func (m *ProjectLogMiddleware) GetProjectByID(ctx context.Context, id string) (project *types.Project, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get project")
		} else {
			logrus.WithFields(logrus.Fields{
				"projectID": project.ID,
				"took":      time.Since(start),
			}).Info("Project retrieved succesfully")
		}
	}(time.Now())
	project, err = m.next.GetProjectByID(ctx, id)
	return project, err

}
func (m *ProjectLogMiddleware) AddTask(ctx context.Context, projectID string, params types.AddTaskParams) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to add task to project")
		} else {
			logrus.WithFields(logrus.Fields{
				"projectID": projectID,
				"taskID":    params.TaskID,
				"took":      time.Since(start),
			}).Info("Task added to project succesfully")
		}
	}(time.Now())
	err = m.next.AddTask(ctx, projectID, params)
	return err
}
