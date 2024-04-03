package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params data.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrInvalidCredentials()
		}
		return err
	}
	if !user.IsPasswordValid(params.Password) {
		return ErrInvalidCredentials()
	}
	if !user.Enabled {
		return ErrForbidden()
	}
	fmt.Println("authenticated -> ", user)
	token := CreateTokenFromUser(user)
	resp := data.AuthResponse{
		Token: token,
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *data.User) string {
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 4).Unix(),
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to signed token with secret", err)
	}
	return tokenStr
}
