package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

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
