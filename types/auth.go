package types

import (
	"fmt"
	"time"

	"github.com/twinj/uuid"
)

type Auth struct {
	UserID         string `bson:"userID" dynamodbav:"userID"`
	AuthUUID       string `bson:"authUUID" dynamodbav:"authUUID"`
	ExpirationTime int64  `bson:"expirationTime" dynamodbav:"expirationTime"`
}

func NewAuth(userID string) *Auth {
	return &Auth{
		UserID:         userID,
		AuthUUID:       uuid.NewV4().String(),
		ExpirationTime: time.Now().Add(time.Hour * 4).Unix(),
	}
}

type AuthFilter struct {
	UserID   string `dynamodbav:"userID"`
	AuthUUID string `dynamodbav:"authUUID"`
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p AuthParams) Validate() map[string]string {
	errors := map[string]string{}
	if !isEmailValid(p.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", p.Email)
	}
	if len(p.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length must be at least %d characters", minPasswordLen)
	}
	return errors
}

type AuthResponse struct {
	Token string `json:"token"`
}
