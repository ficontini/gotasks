package service

import "github.com/ficontini/gotasks/db"

type Service struct {
	Auth    *AuthService
	User    *UserService
	Task    *TaskService
	Project ProjectServicer
}

func NewService(store *db.Store) *Service {
	return &Service{
		Auth:    NewAuthService(store),
		User:    NewUserService(store),
		Task:    NewTaskService(store),
		Project: NewLogMiddleware(NewProjectService(store)),
	}
}
