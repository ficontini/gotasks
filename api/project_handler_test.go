package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
)

func TestPostProjectSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		apiv1          = app.Group("/", JWTAuthentication(db.User))
		projectService = service.NewProjectService(*db.Store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(db.Store, "james", "foo", "supersecure", false, true)
	)
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := data.CreateProjectParams{
		Title:       "test-project",
		Description: "description of this project",
	}
	res, err := makePostRequest(params, user, app)
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
		apiv1          = app.Group("/", JWTAuthentication(db.User))
		projectService = service.NewProjectService(*db.Store)
		projectHandler = NewProjectHandler(projectService)
		user           = fixtures.AddUser(db.Store, "james", "foo", "supersecure", false, true)
	)
	apiv1.Post("/", projectHandler.HandlePostProject)
	params := data.CreateProjectParams{
		Title:       "",
		Description: "description of this project",
	}
	res, err := makePostRequest(params, user, app)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusBadRequest, res.StatusCode)
}

func makePostRequest(params interface{}, user *data.User, app *fiber.App) (*http.Response, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", CreateTokenFromUser(user)))
	req.Header.Add("Content-Type", "application/json")

	return app.Test(req)
}
func decodeToProject(response *http.Response) *data.Project {
	var project *data.Project
	json.NewDecoder(response.Body).Decode(&project)
	return project
}
