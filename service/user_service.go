package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

type UserService struct {
	userStore db.UserStore
}

func NewUserService(userStore db.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (svc *UserService) CreateUser(ctx context.Context, params types.CreateUserParams) (*types.User, error) {
	if svc.isEmailAlreadyInUse(ctx, params.Email) {
		return nil, ErrEmailAlreadyInUse
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return nil, err
	}
	return svc.userStore.InsertUser(ctx, user)

}

func (svc *UserService) isEmailAlreadyInUse(ctx context.Context, email string) bool {
	user, _ := svc.userStore.GetUserByEmail(ctx, email)
	return user != nil
}

func (svc *UserService) setEnabled(ctx context.Context, id string, enabled bool) error {
	user, err := svc.userStore.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrUserNotFound
		}
	}
	if user.Enabled == enabled {
		return ErrUserStateUnchanged
	}
	if err := svc.userStore.Update(ctx, id, types.StatusUpdater{Enabled: enabled}); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (svc *UserService) EnableUser(ctx context.Context, id string) error {
	return svc.setEnabled(ctx, id, true)
}

func (svc *UserService) DisableUser(ctx context.Context, id string) error {
	return svc.setEnabled(ctx, id, false)
}

func (svc *UserService) ResetPassword(ctx context.Context, user *types.User, params types.ResetPasswordParams) error {
	if !user.IsPasswordValid(params.CurrentPassword) {
		return ErrCurrentPassword
	}
	enpw, err := params.GeneratePassword()
	if err != nil {
		return err
	}
	if err := svc.userStore.Update(ctx, user.ID, types.PasswordUpdater{EncryptedPassword: enpw}); err != nil {
		return err
	}
	//TODO: Invalidating existing token
	return nil
}

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrUserStateUnchanged = errors.New("user state unchanged")
	ErrUserNotFound       = errors.New("user resource not found")
	ErrCurrentPassword    = errors.New("current password is not valid")
)
