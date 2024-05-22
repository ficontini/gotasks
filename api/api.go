package api

import (
	"reflect"

	"github.com/ficontini/gotasks/service"
)

type ResourceResponse struct {
	Data    any   `json:"data"`
	Results int   `json:"results"`
	Page    int64 `json:"page"`
}

func NewResourceResponse(data any, results int, page int64) ResourceResponse {
	//TODO: Review
	if reflect.ValueOf(data).IsNil() {
		data = []any{}
	}
	return ResourceResponse{
		Data:    data,
		Results: results,
		Page:    page,
	}
}

type Handler struct {
	Auth    *AuthHandler
	User    *UserHandler
	Task    *TaskHandler
	Project *ProjectHandler
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		Auth:    NewAuthHandler(svc.Auth),
		User:    NewUserHandler(svc.User),
		Task:    NewTaskHandler(svc.Task),
		Project: NewProjectHandler(svc.Project),
	}
}
