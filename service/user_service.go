package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
)

type UserService struct {
	userStore db.UserStore
}

func NewUserService(userStore db.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (svc *UserService) CreateUser(ctx context.Context, params data.CreateUserParams) (*data.User, error) {
	if svc.isEmailAlreadyInUse(ctx, params.Email) {
		return nil, ErrEmailAlreadyInUse
	}
	user, err := data.NewUserFromParams(params)
	if err != nil {
		return nil, err
	}
	return svc.userStore.InsertUser(ctx, user)

}

func (svc *UserService) isEmailAlreadyInUse(ctx context.Context, email string) bool {
	user, _ := svc.userStore.GetUserByEmail(ctx, email)
	return user != nil
}

func (svc *UserService) EnableUser(ctx context.Context, id string) error {
	user, err := svc.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Enabled {
		return ErrConflict
	}
	if err := svc.userStore.Update(ctx, id, db.Map{"enabled": true}); err != nil {
		return err
	}
	return nil
}

var (
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrConflict          = errors.New("user is already enabled")
)
