package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		password    = "supersecurepassword"
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", password, false, true)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		authHandler = NewAuthHandler(authService)
		params      = types.AuthParams{
			Email:    user.Email,
			Password: password,
		}
	)
	app.Post("/auth", authHandler.HandleAuthenticate)
	b := marshallParamsToJSON(t, params)
	req := makeUnauthenticatedRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	resp := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, resp.StatusCode)
	var authresp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authresp); err != nil {
		log.Fatal(err)
	}
	if len(authresp.Token) == 0 {
		log.Fatal("expected the JWT token to be present in the auth response")
	}

}
func TestAuthenticateWrongWithPasswordFailure(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", "supersecurepassword", false, true)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		authHandler = NewAuthHandler(authService)
		params      = types.AuthParams{
			Email:    user.Email,
			Password: "wrongpassword",
		}
	)
	app.Post("/auth", authHandler.HandleAuthenticate)
	b := marshallParamsToJSON(t, params)
	req := makeUnauthenticatedRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	resp := testRequest(t, app, req)
	checkStatusCode(t, http.StatusUnauthorized, resp.StatusCode)
}
