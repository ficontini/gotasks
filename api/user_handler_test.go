package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
)

func TestEnableUserSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		password  = "supersecurepassword"
		user      = fixtures.AddUser(db.Store, "james", "foo", password, false, false)
		adminUser = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		app       = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		apiv1     = app.Group("/", JWTAuthentication(db.User))
		admin     = apiv1.Group("/admin", AdminAuth)
		handler   = NewUserHandler(*service.NewUserService(db.User))
	)
	admin.Post("/user/:id/enable", handler.HandleEnableUser)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/user/%s/enable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", CreateTokenFromUser(adminUser)))
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
		password  = "supersecurepassword"
		user      = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		adminUser = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		app       = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		apiv1     = app.Group("/", JWTAuthentication(db.User))
		admin     = apiv1.Group("/admin", AdminAuth)
		handler   = NewUserHandler(*service.NewUserService(db.User))
	)
	admin.Post("/user/:id/enable", handler.HandleEnableUser)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/user/%s/enable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", CreateTokenFromUser(adminUser)))
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
		password  = "supersecurepassword"
		user      = fixtures.AddUser(db.Store, "james", "foo", password, false, true)
		adminUser = fixtures.AddUser(db.Store, "admin", "foo", password, true, true)
		app       = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		apiv1     = app.Group("/", JWTAuthentication(db.User))
		admin     = apiv1.Group("/admin", AdminAuth)
		handler   = NewUserHandler(*service.NewUserService(db.User))
	)
	admin.Post("/user/:id/disable", handler.HandleDisableUser)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/user/%s/disable", user.ID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", CreateTokenFromUser(adminUser)))
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
		adminUser = fixtures.AddUser(db.Store, "admin", "foo", "supersecurepassword", true, true)
		app       = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		apiv1     = app.Group("/", JWTAuthentication(db.User))
		admin     = apiv1.Group("/admin", AdminAuth)
		handler   = NewUserHandler(*service.NewUserService(db.User))
		wrongID   = "609c4b22a2c2d9c3f83a01f6"
	)
	admin.Post("/user/:id/disable", handler.HandleDisableUser)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/user/%s/disable", wrongID), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", CreateTokenFromUser(adminUser)))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	checkStatusCode(t, http.StatusNotFound, res.StatusCode)
}
