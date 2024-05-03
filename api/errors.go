package api

import (
	"net/http"

	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var (
		fields = logrus.Fields{}
	)
	user, ok := c.Context().Value("user").(*types.User)
	if ok {
		fields["userID"] = user.ID
	}
	if apiError, ok := err.(Error); ok {
		fields["err"] = apiError
		logrus.WithFields(fields).Error()
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	fields["err"] = apiError
	logrus.WithFields(fields).Error()
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}
func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}
func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}
func ErrBadRequestCustomMessage(msg string) Error {
	return NewError(http.StatusBadRequest, msg)
}
func ErrResourceNotFound(msg string) Error {
	return NewError(http.StatusNotFound, msg)
}
func ErrUnAuthorized() Error {
	return NewError(http.StatusUnauthorized, "unathorized request")
}
func ErrInvalidCredentials() Error {
	return NewError(http.StatusUnauthorized, "invalid credentials")
}
func ErrForbidden() Error {
	return NewError(http.StatusForbidden, "user account is not enabled. Please contact the administrator for assistance.")
}
func ErrConflict(msg string) Error {
	return NewError(http.StatusConflict, msg)
}
func ErrInternalServer() Error {
	return NewError(http.StatusInternalServerError, "internal server error")
}
