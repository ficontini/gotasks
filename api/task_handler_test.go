package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
	)

	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
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
	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
	)

	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
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

	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
	)
	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
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

	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
	)
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
		taskService  = service.NewTaskService(db.Store)
		taskHandler  = NewTaskHandler(taskService)
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
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		wrongID     = "609c4b22a2c2d9c3f83a01f6"
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
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(db.Store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, james.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
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
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
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
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
		auth        = fixtures.AddAuth(db.Store, james.ID)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, james.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestCompleteTaskWithAnotherAssignedUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		auth        = fixtures.AddAuth(db.Store, james.ID)
		tom         = fixtures.AddUser(db.Store, "tom", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, tom.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)

}
func TestUpdateDueDateTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(db.Store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, james.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, 1, 5),
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
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
	updatedTask, err := db.Store.Task.GetTaskByID(context.Background(), task.ID)
	if err != nil {
		log.Fatal(err)
	}
	if updatedTask.DueDate != params.DueDate {
		t.Fatalf("expected %v, but got %v", params.DueDate, updatedTask.DueDate)
	}

}
func TestUpdateDueDateTaskWithWrongDate(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(db.Store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, james.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, -3, 5),
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestUpdateDueDateTaskWithAnotherAssignedUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(db.Store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false, true)
		auth        = fixtures.AddAuth(db.Store, james.ID)
		tom         = fixtures.AddUser(db.Store, "tom", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(db.Store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
	)
	fixtures.AssignTaskToUser(*db.Store, task.ID, tom.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, 1, 5),
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)

}
func sendPostRequest(t *testing.T, app *fiber.App, params types.NewTaskParams, expectedStatus int) *types.Task {
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
