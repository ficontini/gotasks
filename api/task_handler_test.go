package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func sendPostRequest(t *testing.T, app *fiber.App, params types.CreateTaskParams, expectedStatus int) *types.Task {
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != expectedStatus {
		t.Fatalf("expected %d status code but got %d", expectedStatus, res.StatusCode)
	}
	var task *types.Task
	json.NewDecoder(res.Body).Decode(&task)
	return task
}

func TestPostTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New()
		taskHandler = NewTaskHandler(db.Task)
	)

	app.Post("/", taskHandler.HandlePostTask)

	params := types.CreateTaskParams{
		Title:       "fake-task",
		Description: "fake description",
		DueDate:     time.Now().AddDate(0, 0, 5),
	}
	task := sendPostRequest(t, app, params, http.StatusOK)
	if len(task.ID) == 0 {
		t.Fatalf("expecting a task id to be set")
	}
	if task.Title != params.Title {
		t.Fatalf("expected title %s but got %s", params.Title, task.Title)
	}
	if task.Description != params.Description {
		t.Fatalf("expected description %s but got %s", params.Description, task.Description)
	}
	if task.Completed {
		t.Fatalf("expected task not completed")
	}
}
func TestPostTaskWithWrongDueDate(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	taskHandler := NewTaskHandler(db.Task)
	app.Post("/", taskHandler.HandlePostTask)

	params := types.CreateTaskParams{
		Title:       "fake-task",
		Description: "fake description",
		DueDate:     time.Now().AddDate(-1, 0, 5),
	}
	task := sendPostRequest(t, app, params, http.StatusBadRequest)
	if len(task.ID) > 0 {
		t.Fatalf("task shouldn't be created")
	}
}
func TestPostInvalidTitle(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	taskHandler := NewTaskHandler(db.Task)
	app.Post("/", taskHandler.HandlePostTask)

	params := types.CreateTaskParams{
		Title:       "aa",
		Description: "fake description",
		DueDate:     time.Now().AddDate(0, 0, 5),
	}
	task := sendPostRequest(t, app, params, http.StatusBadRequest)
	if len(task.ID) > 0 {
		t.Fatalf("task shouldn't be created")
	}
}

func TestPostEmptyRequestBody(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	taskHandler := NewTaskHandler(db.Task)
	app.Post("/", taskHandler.HandlePostTask)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}"))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d status code but got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestDeleteTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		insertedTask = fixtures.AddTask(db.Store, "fake-task", "fake task description", time.Now().AddDate(0, 0, 2), false)
		app          = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskHandler  = NewTaskHandler(db.Task)
	)
	app.Delete("/:id", taskHandler.HandleDeleteTask)
	app.Get("/:id", taskHandler.HandleGetTask)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", insertedTask.ID), nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected %d status code, but got %d", http.StatusOK, res.StatusCode)
	}
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", insertedTask.ID), nil)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %d status code, but got %d", http.StatusNotFound, res.StatusCode)
	}
}
