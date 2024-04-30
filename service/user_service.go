package service

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
)

type UserInserter interface {
	CreateUser(context.Context, types.CreateUserParams) (*types.User, error)
}
type UserGetter interface {
	GetUsers(context.Context, UserQueryParams) ([]*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
}
type UserUpdater interface {
	EnableUser(context.Context, string) error
	DisableUser(context.Context, string) error
	ResetPassword(context.Context, *types.User, types.ResetPasswordParams) error
	InvalidateJWT(context.Context, *types.Auth) error
}
type UserServicer interface {
	UserInserter
	UserGetter
	UserUpdater
}
type UserService struct {
	store *db.Store
}

func NewUserService(store *db.Store) UserServicer {
	return &UserService{
		store: store,
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
	return svc.store.User.InsertUser(ctx, user)

}

func (svc *UserService) isEmailAlreadyInUse(ctx context.Context, email string) bool {
	user, _ := svc.store.User.GetUserByEmail(ctx, email)
	return user != nil
}

func (svc *UserService) setEnabled(ctx context.Context, id string, enabled bool) error {
	user, err := svc.store.User.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrUserNotFound
		}
	}
	if user.Enabled == enabled {
		return ErrUserStateUnchanged
	}
	if err := svc.store.User.Update(ctx, id, db.StatusUpdater{Enabled: enabled}); err != nil {
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
	if err := svc.store.User.Update(ctx, user.ID, db.PasswordUpdater{EncryptedPassword: enpw}); err != nil {
		return err
	}
	return nil
}
func (svc *UserService) InvalidateJWT(ctx context.Context, auth *types.Auth) error {
	filter := &types.AuthFilter{
		UserID:   auth.UserID,
		AuthUUID: auth.AuthUUID,
	}
	if err := svc.store.Auth.Delete(ctx, filter); err != nil {
		return err
	}
	return nil
}

type UserQueryParams struct {
	db.Pagination
}

func (svc *UserService) GetUsers(ctx context.Context, params UserQueryParams) ([]*types.User, error) {
	filter := db.EmptyFilter{DataType: &db.DataType{DataType: types.UserDataType}}
	return svc.store.User.GetUsers(ctx, filter, &params.Pagination)
}
func (svc *UserService) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	user, err := svc.store.User.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrUserStateUnchanged = errors.New("user state unchanged")
	ErrUserNotFound       = errors.New("user resource not found")
	ErrCurrentPassword    = errors.New("current password is not valid")
)
