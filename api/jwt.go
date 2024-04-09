package api

import (
	"fmt"

	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("--- JWT authentication")

		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			fmt.Println("token is not present in the header")
			return ErrUnAuthorized()
		}
		tokenStr := token[0]
		claims, err := authService.ValidateToken(tokenStr[len("Bearer "):])
		if err != nil {
			return ErrUnAuthorized()
		}
		user, err := authService.GetUser(c.Context(), claims)
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		auth, err := authService.GetAuth(c.Context(), claims)
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("auth", auth)
		fmt.Println("token", token)
		return c.Next()
	}
}
