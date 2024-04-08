package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
)

func AddProject(store *db.Store, title, description string, userID string, tasks []string) *data.Project {
	project := &data.Project{
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
func AddTask(store *db.Store, title, description string, dueTo time.Time, completed bool) *data.Task {
	task := &data.Task{
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
func AddUser(store *db.Store, fn, ln, pwd string, isAdmin, enabled bool) *data.User {
	user, err := data.NewUserFromParams(data.CreateUserParams{
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
func AssignTaskToUser(store db.Store, taskID, userID string) {
	if err := store.Task.SetTaskAssignee(context.Background(), taskID, data.TaskAssignmentRequest{UserID: userID}); err != nil {
		log.Fatal(err)
	}

}
