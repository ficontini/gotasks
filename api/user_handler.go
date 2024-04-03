package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params data.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	if h.isEmailAlreadyInUse(c.Context(), params.Email) {
		return ErrBadRequestCustomMessage("email already in use")
	}
	user, err := data.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
func (h *UserHandler) HandleEnableUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	if err := h.userStore.Update(c.Context(), data.ID(id), db.Map{"enabled": true}); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("user")
		}
		return err
	}
	return c.JSON(fiber.Map{"enabled": id})
}
func (h *UserHandler) isEmailAlreadyInUse(ctx context.Context, email string) bool {
	user, _ := h.userStore.GetUserByEmail(ctx, email)
	return user != nil
}
