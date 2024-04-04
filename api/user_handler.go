package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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
	insertedUser, err := h.userService.CreateUser(c.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyInUse) {
			return ErrBadRequestCustomMessage(err.Error())
		}
		return err
	}
	return c.JSON(insertedUser)
}
func (h *UserHandler) HandleEnableUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	if err := h.userService.EnableUser(c.Context(), id); err != nil {
		switch {
		case errors.Is(err, db.ErrorNotFound):
			return ErrResourceNotFound("user")
		case errors.Is(err, service.ErrConflict):
			return ErrConflict(err.Error())
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{"enabled": id})
}
