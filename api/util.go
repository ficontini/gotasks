package api

import (
	"github.com/ficontini/gotasks/data"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*data.User, error) {
	user, ok := c.Context().Value("user").(*data.User)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return user, nil
}
