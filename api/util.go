package api

import (
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return user, nil
}
