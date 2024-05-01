package api

import (
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		logrus.Error("unauthorized user")
		return ErrUnAuthorized()
	}
	if !user.IsAdmin {
		logrus.Error("unauthorized user")
		return ErrUnAuthorized()
	}
	return c.Next()
}
