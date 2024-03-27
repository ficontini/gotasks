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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
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
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", insertedTask.ID), nil)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestDeleteTaskWithWrongID(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskHandler = NewTaskHandler(db.Task)
		wrongID, _  = primitive.ObjectIDFromHex("notavalidobjectid")
	)
	app.Delete("/:id", taskHandler.HandleDeleteTask)
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", wrongID), nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestCompleteTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskHandler = NewTaskHandler(db.Task)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
	)
	app.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	var result map[string]string
	json.NewDecoder(res.Body).Decode(&result)
	if result["updated"] != task.ID {
		t.Fatal("updating a different task")
	}
	app.Get("/:id", taskHandler.HandleGetTask)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", task.ID), nil)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var updatedTask *types.Task
	json.NewDecoder(res.Body).Decode(&updatedTask)
	if !updatedTask.Completed {
		t.Fatalf("task wiht %s expected complete", updatedTask.ID)
	}
}
func TestCompleteTaskWithCompletedStatus(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskHandler = NewTaskHandler(db.Task)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
	)
	app.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func sendPostRequest(t *testing.T, app *fiber.App, params types.CreateTaskParams, expectedStatus int) *types.Task {
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, expectedStatus, res.StatusCode)
	var task *types.Task
	json.NewDecoder(res.Body).Decode(&task)
	return task
}

func checkStatusCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("expected %d status code, but got %d", expected, actual)
	}
}