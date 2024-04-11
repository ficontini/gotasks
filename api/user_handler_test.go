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

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func TestEnableUserSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepassword"
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false, false)
		adminUser   = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(db.Store))
		auth        = fixtures.AddAuth(db.Store, adminUser.ID)
	)
	admin.Put("/user/:id/enable", handler.HandleEnableUser)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/enable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	updatedUser, err := db.Store.User.GetUserByID(context.Background(), user.ID)
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
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		adminUser   = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		auth        = fixtures.AddAuth(db.Store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(db.Store))
	)
	admin.Put("/user/:id/enable", handler.HandleEnableUser)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/enable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusConflict, res.StatusCode)
}
func TestDisableUserSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepassword"
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		adminUser   = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		auth        = fixtures.AddAuth(db.Store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(db.Store))
	)
	admin.Put("/user/:id/disable", handler.HandleDisableUser)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/disable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)
	updatedUser, err := db.Store.User.GetUserByID(context.Background(), user.ID)
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
		adminUser   = fixtures.AddUser(db.Store, "admin", "foo", "supersecurepassword", true, true)
		auth        = fixtures.AddAuth(db.Store, adminUser.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		admin       = apiv1.Group("/admin", AdminAuth)
		handler     = NewUserHandler(service.NewUserService(db.Store))
		wrongID     = "609c4b22a2c2d9c3f83a01f6"
	)
	admin.Put("/user/:id/disable", handler.HandleDisableUser)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/admin/user/%s/disable", wrongID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authService.CreateTokenFromAuth(auth)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
func TestResetPasswordSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepwd"
		newpassword = "newsupersecurepwd"
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		auth        = fixtures.AddAuth(db.Store, user.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		handler     = NewUserHandler(service.NewUserService(db.Store))
		token       = authService.CreateTokenFromAuth(auth)
	)
	apiv1.Post("/reset-password", handler.HandleResetPassword)
	params := types.ResetPasswordParams{
		CurrentPassword: password,
		NewPassword:     newpassword,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	res, err := app.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)

	apiv1.Get("/user", handler.HandleGetUser)
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	res, err = app.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)
}
func TestResetPasswordWithWrongCurrentPassword(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password    = "supersecurepwd"
		newpassword = "newsupersecurepwd"
		user        = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		auth        = fixtures.AddAuth(db.Store, user.ID)
		app         = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		authService = service.NewAuthService(db.Store)
		apiv1       = app.Group("/", JWTAuthentication(authService))
		handler     = NewUserHandler(service.NewUserService(db.Store))
		token       = authService.CreateTokenFromAuth(auth)
	)
	apiv1.Post("/reset-password", handler.HandleResetPassword)
	params := types.ResetPasswordParams{
		CurrentPassword: newpassword,
		NewPassword:     newpassword,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	res, err := app.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	checkStatusCode(t, http.StatusUnauthorized, res.StatusCode)
	apiv1.Get("/user", handler.HandleGetUser)
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	res, err = app.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	checkStatusCode(t, http.StatusOK, res.StatusCode)
}
