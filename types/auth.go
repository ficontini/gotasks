package types

import "fmt"

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
