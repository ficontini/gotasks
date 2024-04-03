package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db/fixtures"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		password    = "supersecurepassword"
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authHandler = NewAuthHandler(db.User)
		params      = data.AuthParams{
			Email:    user.Email,
			Password: password,
		}
	)
	app.Post("/auth", authHandler.HandleAuthenticate)
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, resp.StatusCode)
	var authresp data.AuthResponse
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
		user        = fixtures.AddUser(db.Store, "james", "foo", "supersecurepassword", false)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authHandler = NewAuthHandler(db.User)
		params      = data.AuthParams{
			Email:    user.Email,
			Password: "wrongpassword",
		}
	)
	app.Post("/auth", authHandler.HandleAuthenticate)
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusUnauthorized, resp.StatusCode)
}
