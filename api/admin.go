package api

import (
	"github.com/ficontini/gotasks/data"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*data.User)
	if !ok {
		return ErrUnAuthorized()
	}
	if !user.IsAdmin {
		return ErrUnAuthorized()
	}
	return c.Next()
}
