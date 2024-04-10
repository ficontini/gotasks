package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: *userService,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
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
		case errors.Is(err, service.ErrUserNotFound):
			return ErrResourceNotFound(err.Error())
		case errors.Is(err, service.ErrUserStateUnchanged):
			return ErrConflict(err.Error())
		default:
			return err
		}
	}
	return c.JSON(fiber.Map{"enabled": id})
}
func (h *UserHandler) HandleDisableUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	if err := h.userService.DisableUser(c.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return ErrResourceNotFound(err.Error())
		case errors.Is(err, service.ErrUserStateUnchanged):
			return ErrConflict(err.Error())
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{"disabled": id})
}
func (h *UserHandler) HandleResetPassword(c *fiber.Ctx) error {
	user, err := getUserAuth(c)
	if err != nil {
		return err
	}
	auth, err := getAuth(c)
	if err != nil {
		return err
	}
	var params types.ResetPasswordParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	if err := h.userService.ResetPassword(c.Context(), user, params); err != nil {
		if errors.Is(err, service.ErrCurrentPassword) {
			return ErrUnAuthorized()
		}
		return err
	}
	if err := h.userService.InvalidateJWT(c.Context(), auth); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"password": "updated"})
}
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := getUserAuth(c)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
