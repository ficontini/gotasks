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
		if err != nil {
			logrus.WithError(err).Error("Failed to authenticate user")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": auth.UserID,
			}).Info("AuthenticatedUser successfully completed")
		}
	}(time.Now())
	auth, err = m.next.AuthenticateUser(ctx, params)
	return auth, err
}

func (m *AuthLogMiddleware) GetUser(ctx context.Context, claims jwt.MapClaims) (user *types.User, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get user")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": user.ID,
			}).Info("Get user")
		}
	}(time.Now())
	user, err = m.next.GetUser(ctx, claims)
	return user, err
}

func (m *AuthLogMiddleware) GetAuth(ctx context.Context, claims jwt.MapClaims) (auth *types.Auth, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to get auth")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start),
				"userID": auth.UserID,
			}).Info("Get auth")
		}
	}(time.Now())
	auth, err = m.next.GetAuth(ctx, claims)
	return auth, err
}

func (m *AuthLogMiddleware) CreateTokenFromAuth(auth *types.Auth) (token string, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to create token from auth")
		}
	}(time.Now())
	token, err = m.next.CreateTokenFromAuth(auth)
	return token, err
}

func (m *AuthLogMiddleware) ValidateToken(tokenStr string) (claims jwt.MapClaims, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithError(err).Error("Failed to validate token")
		}
	}(time.Now())
	claims, err = m.next.ValidateToken(tokenStr)
	return claims, err
}
