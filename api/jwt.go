package api

import (
	"fmt"
	"os"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("--- JWT authentication")

		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			fmt.Println("token is not present in the header")
			return ErrUnAuthorized()
		}
		tokenStr := token[0]
		claims, err := validateToken(tokenStr[len("Bearer "):])
		if err != nil {
			return ErrUnAuthorized()
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), types.ID(userID))
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		fmt.Println("token", token)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return claims, nil
}
