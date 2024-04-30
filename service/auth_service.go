package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
	"github.com/golang-jwt/jwt/v5"
)

type AuthGetter interface {
	GetUser(context.Context, jwt.MapClaims) (*types.User, error)
	GetAuth(context.Context, jwt.MapClaims) (*types.Auth, error)
}
type AuthServicer interface {
	AuthGetter
	AuthenticateUser(context.Context, *types.AuthParams) (*types.Auth, error)
	CreateTokenFromAuth(*types.Auth) (string, error)
	ValidateToken(string) (jwt.MapClaims, error)
}
type AuthService struct {
	store *db.Store
}

func NewAuthService(store *db.Store) AuthServicer {
	return &AuthService{
		store: store,
	}
}
func (svc *AuthService) AuthenticateUser(ctx context.Context, params *types.AuthParams) (*types.Auth, error) {
	user, err := svc.store.User.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !user.IsPasswordValid(params.Password) {
		return nil, ErrInvalidCredentials
	}
	if !user.Enabled {
		return nil, ErrForbidden
	}
	auth, err := svc.store.Auth.Insert(ctx, types.NewAuth(user.ID))
	if err != nil {
		return nil, err
	}
	return auth, nil
}
func (svc *AuthService) GetUser(ctx context.Context, claims jwt.MapClaims) (*types.User, error) {
	user, err := svc.store.User.GetUserByID(ctx, claims["id"].(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (svc *AuthService) GetAuth(ctx context.Context, claims jwt.MapClaims) (*types.Auth, error) {

	filter := &types.AuthFilter{
		UserID:   claims["id"].(string),
		AuthUUID: claims["auth_uuid"].(string),
	}
	auth, err := svc.store.Auth.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
func (svc *AuthService) CreateTokenFromAuth(auth *types.Auth) (string, error) {
	claims := jwt.MapClaims{
		"id":        auth.UserID,
		"auth_uuid": auth.AuthUUID,
		"exp":       auth.ExpirationTime,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return tokenStr, fmt.Errorf("failed to generate auth token")
	}
	return tokenStr, nil
}
func (svc *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnAuthorized
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrUnAuthorized
	}
	if !token.Valid {
		return nil, ErrUnAuthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized
	}
	return claims, nil
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
)
