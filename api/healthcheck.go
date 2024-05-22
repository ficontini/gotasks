package api

import "github.com/gofiber/fiber/v2"

func HandleHealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"healthy": true})
}
