package api

import (
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

func getAuth(c *fiber.Ctx) (*types.Auth, error) {
	auth, ok := c.Context().Value("auth").(*types.Auth)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return auth, nil
}

func getUserAuth(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return user, nil
}
