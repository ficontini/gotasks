package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func TestEnableUserSuccess(t *testing.T) {
	db := setup(t)

	defer db.teardown(t)
	var (
		store       = db.Store()
		password    = "supersecurepassword"
		user        = fixtures.AddUser(store, "james", "foo", password, false, false)
		adminUser   = fixtures.AddUser(store, "admin", "foo", password, true, true)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(store))
		auth        = fixtures.AddAuth(store, adminUser.ID)
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	admin.Put("/user/:id/enable", handler.HandleEnableUser)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/enable", user.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	updatedUser, err := store.User.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !updatedUser.Enabled {
		t.Fatal("expected user to be  enabled")
	}
}
func TestEnableUserAlreadyEnabled(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepassword"
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", password, false, true)
		adminUser   = fixtures.AddUser(store, "admin", "foo", password, true, true)
		auth        = fixtures.AddAuth(store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(store))
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	admin.Put("/user/:id/enable", handler.HandleEnableUser)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/enable", user.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusConflict, res.StatusCode)

}
func TestDisableUserSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepassword"
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", password, false, true)
		adminUser   = fixtures.AddUser(store, "admin", "foo", password, true, true)
		auth        = fixtures.AddAuth(store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(store))
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	admin.Put("/user/:id/disable", handler.HandleDisableUser)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/disable", user.ID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	updatedUser, err := store.User.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedUser.Enabled {
		t.Fatal("expected user not to be enabled")
	}
}
func TestDisableUserWithWrongID(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		store       = db.Store()
		adminUser   = fixtures.AddUser(store, "admin", "foo", "supersecurepassword", true, true)
		auth        = fixtures.AddAuth(store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(store))
		wrongID     = "609c4b22a2c2d9c3f83a01f6"
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	admin.Put("/user/:id/disable", handler.HandleDisableUser)
	req := makeRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/disable", wrongID), token, nil)
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestResetPasswordSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepwd"
		newpassword = "newsupersecurepwd"
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", password, false, true)
		auth        = fixtures.AddAuth(store, user.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		handler     = NewUserHandler(service.NewUserService(store))
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	apiv1.Post("/reset-password", handler.HandleResetPassword)
	params := types.ResetPasswordParams{
		CurrentPassword: password,
		NewPassword:     newpassword,
	}
	b, _ := json.Marshal(params)
	req := makeRequest(http.MethodPost, "/reset-password", token, bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)

	apiv1.Get("/user", handler.HandleGetUser)
	req = makeRequest(http.MethodGet, "/user", token, nil)
	res = testRequest(t, app, req)
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)
}
func TestResetPasswordWithWrongCurrentPassword(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepwd"
		newpassword = "newsupersecurepwd"
		store       = db.Store()
		user        = fixtures.AddUser(store, "james", "foo", password, false, true)
		auth        = fixtures.AddAuth(store, user.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		handler     = NewUserHandler(service.NewUserService(store))
	)
	token, err := authService.CreateTokenFromAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	apiv1.Post("/reset-password", handler.HandleResetPassword)
	params := types.ResetPasswordParams{
		CurrentPassword: newpassword,
		NewPassword:     newpassword,
	}
	b, _ := json.Marshal(params)
	req := makeRequest(http.MethodPost, "/reset-password", token, bytes.NewReader(b))
	res := testRequest(t, app, req)
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)
	apiv1.Get("/user", handler.HandleGetUser)
	req = makeRequest(http.MethodGet, "/user", token, nil)
	res = testRequest(t, app, req)
	checkStatusCode(t, http.StatusOK, res.StatusCode)
}
