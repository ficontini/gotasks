package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

func AddProject(store *db.Store, name, description string, userID string, tasks []string) *types.Project {
	project := &types.Project{
		Name:        name,
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
func AddTask(store *db.Store, name, description string, dueTo time.Time, completed bool) *types.Task {
	task := types.NewTaskFromParams(types.NewTaskParams{
		Name:        name,
		Description: description,
		DueDate:     dueTo,
	})
	task.Completed = completed
	insertedTask, err := store.Task.InsertTask(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	return insertedTask
}
func AddProjectIDToTask(store *db.Store, task *types.Task, projectID string) {
	if err := store.Task.Update(context.Background(), task.ID, db.TaskProjectIDUpdater{ProjectID: projectID}); err != nil {
		log.Fatal(err)
	}
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
func AssignTaskToUser(store *db.Store, taskID, userID string) {
	if err := store.Task.Update(context.Background(), taskID, db.TaskAssignationUpdater{AssignedTo: userID}); err != nil {
		log.Fatal(err)
	}

}
