package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		taskService = service.NewTaskService(db.Store())
		taskHandler = NewTaskHandler(taskService)
	)

	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
		Title:       "fake-task",
		Description: "fake description",
		DueDate:     time.Now().AddDate(0, 0, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeUnauthenticatedRequest(http.MethodPost, "/", bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	task := decodeToTask(t, res)
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
		taskService = service.NewTaskService(db.Store())
		taskHandler = NewTaskHandler(taskService)
	)

	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
		Title:       "fake-task",
		Description: "fake description",
		DueDate:     time.Now().AddDate(-1, 0, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeUnauthenticatedRequest(http.MethodPost, "/", bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestPostInvalidTitle(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store())
		taskHandler = NewTaskHandler(taskService)
	)
	app.Post("/", taskHandler.HandlePostTask)

	params := types.NewTaskParams{
		Title:       "aa",
		Description: "fake description",
		DueDate:     time.Now().AddDate(0, 0, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeUnauthenticatedRequest(http.MethodPost, "/", bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
	task := decodeToTask(t, res)
	if len(task.ID) > 0 {
		t.Fatalf("task shouldn't be created")
	}
}

func TestPostEmptyRequestBody(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New()
		taskService = service.NewTaskService(db.Store())
		taskHandler = NewTaskHandler(taskService)
	)
	app.Post("/", taskHandler.HandlePostTask)

	req := makeUnauthenticatedRequest(http.MethodPost, "/", bytes.NewBufferString("{}"))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		store        = db.Store()
		insertedTask = fixtures.AddTask(store, "fake-task", "fake task description", time.Now().AddDate(0, 0, 2), false)
		app          = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskService  = service.NewTaskService(store)
		taskHandler  = NewTaskHandler(taskService)
	)
	app.Delete("/:id", taskHandler.HandleDeleteTask)
	app.Get("/:id", taskHandler.HandleGetTask)

	req := makeUnauthenticatedRequest(http.MethodDelete, fmt.Sprintf("/%s", insertedTask.ID), nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	req = makeUnauthenticatedRequest(http.MethodGet, fmt.Sprintf("/%s", insertedTask.ID), nil)
	res = testRequest(t, app, req)
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestDeleteTaskWithWrongID(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		taskService = service.NewTaskService(db.Store())
		taskHandler = NewTaskHandler(taskService)
		wrongID     = "609c4b22a2c2d9c3f83a01f6"
	)
	app.Delete("/:id", taskHandler.HandleDeleteTask)
	req := makeUnauthenticatedRequest(http.MethodDelete, fmt.Sprintf("/%s", wrongID), nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestCompleteTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, james.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := makeRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	var result map[string]string
	json.NewDecoder(res.Body).Decode(&result)
	if result["updated"] != task.ID {
		t.Fatal("updating a different task")
	}
	app.Get("/:id", taskHandler.HandleGetTask)
	req = makeRequest(http.MethodGet, fmt.Sprintf("/%s", task.ID), token, nil)
	res = testRequest(t, app, req)
	updatedTask := decodeToTask(t, res)
	if !updatedTask.Completed {
		t.Fatalf("task wiht %s expected complete", updatedTask.ID)
	}
}
func TestCompleteTaskWithCompletedStatus(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
		auth        = fixtures.AddAuth(store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, james.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := makeRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestCompleteTaskWithAnotherAssignedUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		auth        = fixtures.AddAuth(store, james.ID)
		tom         = fixtures.AddUser(store, "tom", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, tom.ID)
	apiv1.Post("/:id/complete", taskHandler.HandleCompleteTask)
	req := makeRequest(http.MethodPost, fmt.Sprintf("/%s/complete", task.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)

}
func TestUpdateDueDateTaskSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, james.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, 1, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), token, bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	var result map[string]string
	json.NewDecoder(res.Body).Decode(&result)
	if result["updated"] != task.ID {
		t.Fatal("updating a different task")
	}
	updatedTask, err := store.Task.GetTaskByID(context.Background(), task.ID)
	if err != nil {
		log.Fatal(err)
	}
	expecedDate := params.DueDate.UTC().Round(time.Second)
	actualDate := updatedTask.DueDate.UTC().Round(time.Second)
	if !actualDate.Equal(expecedDate) {
		t.Fatalf("expected %v, but got %v", expecedDate, actualDate)
	}

}
func TestUpdateDueDateTaskWithWrongDate(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), false)
		auth        = fixtures.AddAuth(store, james.ID)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, james.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, -3, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), token, bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestUpdateDueDateTaskWithAnotherAssignedUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store       = db.Store()
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		taskService = service.NewTaskService(store)
		taskHandler = NewTaskHandler(taskService)
		james       = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		auth        = fixtures.AddAuth(store, james.ID)
		tom         = fixtures.AddUser(store, "tom", "foo", "supersecurepassword", false, true)
		task        = fixtures.AddTask(store, "fake task", "fake task description", time.Now().AddDate(0, 0, 5), true)
		token       = authService.CreateTokenFromAuth(auth)
	)
	fixtures.AssignTaskToUser(store, task.ID, tom.ID)
	apiv1.Put("/:id/due-date", taskHandler.HandlePutDueDateTask)
	params := types.UpdateDueDateTaskRequest{
		DueDate: time.Now().AddDate(0, 1, 5),
	}
	b := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/%s/due-date", task.ID), token, bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)

}

func decodeToTask(t *testing.T, response *http.Response) *types.Task {
	var task *types.Task
	if err := json.NewDecoder(response.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	return task
}
