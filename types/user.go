package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLen      = 7
	minFirstNameLen     = 5
	minLastNameLen      = 5
	defaultCost     int = 12
)

type User struct {
	ID                ID     `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool   `bson:"isAdmin" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), defaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
func (u *User) IsPasswordValid(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil

}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.Password) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length must be at least %d characters", minFirstNameLen)
	}
	if len(p.Password) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length must be at least %d characters", minLastNameLen)
	}
	if len(p.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length must be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(p.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", p.Email)
	}
	return errors
}
func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
