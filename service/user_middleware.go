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
		if err != nil {
			logrus.WithError(err).Error("Failed to create user")
		} else {
			logrus.WithFields(logrus.Fields{
				"userID": user.ID,
				"took":   time.Since(start),
			}).Info("User created successfully")
		}
	}(time.Now())
	user, err = m.next.CreateUser(ctx, params)
	return user, err
}
func (m *UserLogMiddleware) EnableUser(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to enable user")
		} else {
			logrus.WithFields(logrus.Fields{
				"userID": id,
				"took":   time.Since(start),
			}).Info("EnableUser successfully completed")
		}
	}(time.Now())

	err = m.next.EnableUser(ctx, id)
	return err
}
func (m *UserLogMiddleware) DisableUser(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to disable user")
		} else {
			logrus.WithFields(logrus.Fields{
				"userID": id,
				"took":   time.Since(start),
			}).Info("DisableUser successfully completed")
		}
	}(time.Now())

	err = m.next.DisableUser(ctx, id)
	return err
}
func (m *UserLogMiddleware) ResetPassword(ctx context.Context, user *types.User, params types.ResetPasswordParams) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("ResetPassword failed")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": user.ID,
			}).Info("ResetPassword successfully completed")
		}
	}(time.Now())
	err = m.next.ResetPassword(ctx, user, params)
	return err
}
func (m *UserLogMiddleware) InvalidateJWT(ctx context.Context, auth *types.Auth) (err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"userID": auth.UserID,
				"error":  err,
			}).Error("InvalidateJWT: error deleting auth token from database")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": auth.UserID,
			}).Info("InvalidateJWT successfully completed")
		}
	}(time.Now())

	err = m.next.InvalidateJWT(ctx, auth)
	return err
}
func (m *UserLogMiddleware) GetUsers(ctx context.Context, params UserQueryParams) (users []*types.User, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get users")
		} else {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
			}).Info("Get users")
		}
	}(time.Now())
	users, err = m.next.GetUsers(ctx, params)
	return users, err
}

func (m *UserLogMiddleware) GetUserByID(ctx context.Context, id string) (user *types.User, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get user")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": user.ID,
			}).Info("Get user by ID")
		}
	}(time.Now())
	user, err = m.next.GetUserByID(ctx, id)
	return user, err
}
