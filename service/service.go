package service

import "github.com/ficontini/gotasks/db"

type Service struct {
	Auth    AuthServicer
	User    UserServicer
	Task    TaskServicer
	Project ProjectServicer
}

func NewService(store *db.Store) *Service {
	return &Service{
		Auth:    NewAuthLogMiddleware(NewAuthService(store)),
		User:    NewUserLogMiddleware(NewUserService(store)),
		Task:    NewTaskLogMiddleware(NewTaskService(store)),
		Project: NewProjectLogMiddleware(NewProjectService(store)),
	}
}
