package service

import (
	"context"
	"time"

	"github.com/ficontini/gotasks/types"
	"github.com/sirupsen/logrus"
)

type UserLogMiddleware struct {
	next UserServicer
}

func NewUserLogMiddleware(next UserServicer) UserServicer {
	return &UserLogMiddleware{
		next: next,
	}
}

func (m *UserLogMiddleware) CreateUser(ctx context.Context, params types.CreateUserParams) (user *types.User, err error) {
	defer func(start time.Time) {
		var id string
		if user != nil {
			id = user.ID
		}
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   id,
			"err":  err,
		}).Info("Create user")
	}(time.Now())
	user, err = m.next.CreateUser(ctx, params)
	return user, err
}
func (m *UserLogMiddleware) EnableUser(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   id,
			"err":  err,
		}).Info("Enable user")
	}(time.Now())

	err = m.next.EnableUser(ctx, id)
	return err
}
func (m *UserLogMiddleware) DisableUser(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   id,
			"err":  err,
		}).Info("Disable user")
	}(time.Now())

	err = m.next.DisableUser(ctx, id)
	return err
}
func (m *UserLogMiddleware) ResetPassword(ctx context.Context, user *types.User, params types.ResetPasswordParams) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   user.ID,
			"err":  err,
		}).Info("Reset user password")
	}(time.Now())
	err = m.next.ResetPassword(ctx, user, params)
	return err
}
func (m *UserLogMiddleware) InvalidateJWT(ctx context.Context, auth *types.Auth) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"userID": auth.UserID,
			"err":    err,
		}).Info("Invalidate JWT token")
	}(time.Now())

	err = m.next.InvalidateJWT(ctx, auth)
	return err
}
func (m *UserLogMiddleware) GetUsers(ctx context.Context, params UserQueryParams) (users []*types.User, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Get users")
	}(time.Now())
	users, err = m.next.GetUsers(ctx, params)
	return users, err
}

func (m *UserLogMiddleware) GetUserByID(ctx context.Context, id string) (user *types.User, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   id,
			"err":  err,
		}).Info("Get user by ID")
	}(time.Now())
	user, err = m.next.GetUserByID(ctx, id)
	return user, err
}
