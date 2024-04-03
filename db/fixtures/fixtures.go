package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

func AddProject(store *db.Store, title, description string, userID types.ID, tasks []types.ID) *types.Project {
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
		Projects:    []types.ID{},
	}
	insertedTask, err := store.Task.InsertTask(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	return insertedTask
}
func AddUser(store *db.Store, fn, ln, pwd string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  pwd,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
