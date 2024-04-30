package service

import (
	"context"
	"time"

	"github.com/ficontini/gotasks/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type AuthLogMiddleware struct {
	next AuthServicer
}

func NewAuthLogMiddleware(next AuthServicer) AuthServicer {
	return &AuthLogMiddleware{
		next: next,
	}
}
func (m *AuthLogMiddleware) AuthenticateUser(ctx context.Context, params *types.AuthParams) (auth *types.Auth, err error) {
	defer func(start time.Time) {
		var (
			userID   string
			authUUID string
		)
		if auth != nil {
			userID = auth.UserID
			authUUID = auth.AuthUUID
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"userID":   userID,
			"authUUID": authUUID,
			"err":      err,
		}).Info("Authenticate user")
	}(time.Now())
	auth, err = m.next.AuthenticateUser(ctx, params)
	return auth, err
}

func (m *AuthLogMiddleware) GetUser(ctx context.Context, claims jwt.MapClaims) (user *types.User, err error) {
	defer func(start time.Time) {
		var userID string
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"userID": userID,
			"err":    err,
		}).Info("Get user")
	}(time.Now())
	user, err = m.next.GetUser(ctx, claims)
	return user, err
}

func (m *AuthLogMiddleware) GetAuth(ctx context.Context, claims jwt.MapClaims) (auth *types.Auth, err error) {
	defer func(start time.Time) {
		var (
			userID   string
			authUUID string
		)
		if auth != nil {
			userID = auth.UserID
			authUUID = auth.AuthUUID
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"userID":   userID,
			"authUUID": authUUID,
			"err":      err,
		}).Info("Get auth")
	}(time.Now())
	auth, err = m.next.GetAuth(ctx, claims)
	return auth, err
}

func (m *AuthLogMiddleware) CreateTokenFromAuth(auth *types.Auth) (token string) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"userID": auth.UserID,
			"token":  token,
		}).Info("Create token from auth")
	}(time.Now())
	token = m.next.CreateTokenFromAuth(auth)
	return token
}

func (m *AuthLogMiddleware) ValidateToken(tokenStr string) (claims jwt.MapClaims, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Validate token")
	}(time.Now())
	claims, err = m.next.ValidateToken(tokenStr)
	return claims, err
}
