package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostProjectSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		store          = db.Store()
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService    = service.NewAuthService(store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(store, "james", "foo", "supersecure", false, true)
		auth           = fixtures.AddAuth(store, user.ID)
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := types.NewProjectParams{
		Title:       "test-project",
		Description: "description of this project",
	}
	jsonBytes := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPost, "/", token, bytes.NewReader(jsonBytes))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	project := decodeToProject(t, res)
	if project.UserID != user.ID {
		t.Fatalf("user ID mismatch: %s , %s", project.UserID, user.ID)
	}
	if project.Title != params.Title {
		t.Fatalf("expected %s but got %s", params.Title, project.Title)
	}
	if project.Description != params.Description {
		t.Fatalf("expected %s but got %s", params.Description, project.Description)
	}
}
func TestPostProjectInvalidTitle(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		store          = db.Store()
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService    = service.NewAuthService(store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(store, "james", "foo", "supersecure", false, true)
		auth           = fixtures.AddAuth(store, user.ID)
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := types.NewProjectParams{
		Title:       "",
		Description: "description of this project",
	}
	jsonBytes := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPost, "/", token, bytes.NewReader(jsonBytes))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}

func TestAddTaskToProjectSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store          = db.Store()
		authService    = service.NewAuthService(store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(store, "james", "foo", "supersecure", false, true)
		task           = fixtures.AddTask(store, "task01", "description of task01", time.Now().AddDate(0, 0, 2), false)
		project        = fixtures.AddProject(store, "test-project", "test-project-0001", user.ID, []string{})
		auth           = fixtures.AddAuth(store, user.ID)
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	apiv1.Post("project/:id/task", projectHandler.HandlePostTask)
	params := types.AddTaskParams{
		TaskID: task.ID,
	}
	jsonBytes := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPost, fmt.Sprintf("/project/%s/task", project.ID), token, bytes.NewReader(jsonBytes))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	updatedProject, err := store.Project.GetProjectByID(context.Background(), project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !updatedProject.ContainsTask(task.ID) {
		t.Fatalf("expected the project with id %s contains this task: %s", project.ID, task.ID)
	}
	updatedTask, err := store.Task.GetTaskByID(context.Background(), task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedTask.ProjectID != project.ID {
		t.Fatalf("expected the task %s to be associated with project %s", updatedTask.ID, project.ID)
	}

}
func TestAddTaskToProjectAlreadyAdded(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		store          = db.Store()
		authService    = service.NewAuthService(store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(store, "james", "foo", "supersecure", false, true)
		task           = fixtures.AddTask(store, "task01", "description of task01", time.Now().AddDate(0, 0, 2), false)
		project        = fixtures.AddProject(store, "test-project", "test-project-0001", user.ID, []string{task.ID})
		auth           = fixtures.AddAuth(store, user.ID)
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	fixtures.AddProjectIDToTask(store, task, project.ID)
	apiv1.Post("project/:id/task", projectHandler.HandlePostTask)
	params := types.AddTaskParams{
		TaskID: task.ID,
	}
	jsonBytes := marshallParamsToJSON(t, params)
	req := makeRequest(http.MethodPost, fmt.Sprintf("/project/%s/task", project.ID), token, bytes.NewReader(jsonBytes))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusConflict, res.StatusCode)
}

func decodeToProject(t *testing.T, response *http.Response) *types.Project {
	var project *types.Project
	if err := json.NewDecoder(response.Body).Decode(&project); err != nil {
		t.Fatal(err)
	}
	return project
}
