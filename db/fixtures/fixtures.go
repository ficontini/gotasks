package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

func AddProject(store *db.Store, title, description string, userID string, tasks []string) *types.Project {
	project := &types.Project{
		Title:       title,
		Description: description,
		UserID:      userID,
		Tasks:       tasks,
	}
	insertedProject, err := store.Project.InsertProject(context.Background(), project)
	if err != nil {
		log.Fatal(err)
	}
	return insertedProject
}
func AddTask(store *db.Store, title, description string, dueTo time.Time, completed bool) *types.Task {
	task := &types.Task{
		Title:       title,
		Description: description,
		DueDate:     dueTo,
		Completed:   completed,
	}
	insertedTask, err := store.Task.InsertTask(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	return insertedTask
}
func AddUser(store *db.Store, fn, ln, pwd string, isAdmin, enabled bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  pwd,
	})
	user.IsAdmin = isAdmin
	user.Enabled = enabled
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
func AddAuth(store *db.Store, userID string) *types.Auth {
	auth := types.NewAuth(userID)
	insertedAuth, err := store.Auth.Insert(context.Background(), auth)
	if err != nil {
		log.Fatal(err)
	}
	return insertedAuth
}
func AssignTaskToUser(store db.Store, taskID, userID string) {
	update, err := db.NewTaskAssignationUpdater(userID)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Task.Update(context.Background(), taskID, update); err != nil {
		log.Fatal(err)
	}

}
