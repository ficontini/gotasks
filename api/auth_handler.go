package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthServicer
}

func NewAuthHandler(authService service.AuthServicer) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params types.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	auth, err := h.authService.AuthenticateUser(c.Context(), &params)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return ErrInvalidCredentials()
		case errors.Is(err, service.ErrForbidden):
			return ErrForbidden()

		default:
			return err
		}
	}
	token, err := h.authService.CreateTokenFromAuth(auth)
	if err != nil {
		return err
	}
	resp := &types.AuthResponse{
		Token: token,
	}
	return c.JSON(resp)
}
