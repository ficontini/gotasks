package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostProjectSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService    = service.NewAuthService(db.Store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(db.Store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(db.Store, "james", "foo", "supersecure", false, true)
		auth           = fixtures.AddAuth(db.Store, user.ID)
	)
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := types.NewProjectParams{
		Title:       "test-project",
		Description: "description of this project",
	}
	res, err := makePostRequest(params, authService, auth, app)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	project := decodeToProject(res)
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
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService    = service.NewAuthService(db.Store)
		apiv1          = app.Group("/", JWTAuthentication(authService))
		projectService = service.NewProjectService(db.Store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(db.Store, "james", "foo", "supersecure", false, true)
		auth           = fixtures.AddAuth(db.Store, user.ID)
	)
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := types.NewProjectParams{
		Title:       "",
		Description: "description of this project",
	}
	res, err := makePostRequest(params, authService, auth, app)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}
func TestAddTaskToProjectSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

}

func makePostRequest(params interface{}, authService *service.AuthService, auth *types.Auth, app *fiber.App) (*http.Response, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	req.Header.Add("Content-Type", "application/json")

	return app.Test(req)
}
func decodeToProject(response *http.Response) *types.Project {
	var project *types.Project
	json.NewDecoder(response.Body).Decode(&project)
	return project
}
