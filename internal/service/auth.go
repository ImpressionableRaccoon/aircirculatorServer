package service

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

func (s *Service) SignUp(ctx context.Context, login string, password string) (token string, err error) {
	if login == "" {
		return "", ErrLoginIsEmpty
	}
	if password == "" {
		return "", ErrPasswordIsEmpty
	}

	hash, salt, err := utils.PreparePassword(password, s.cfg.PasswordSalt)
	if err != nil {
		return "", err
	}

	err = s.RegisterUser(ctx, login, hash, salt)
	if err != nil {
		return "", err
	}

	return s.SignIn(ctx, login, password)
}

func (s *Service) SignIn(ctx context.Context, login string, password string) (token string, err error) {
	if login == "" {
		return "", ErrLoginIsEmpty
	}
	if password == "" {
		return "", ErrPasswordIsEmpty
	}

	user, err := s.GetUserByLogin(ctx, login)
	if errors.Is(err, storage.ErrUserNotFound) {
		return "", ErrUnauthorized
	}
	if err != nil {
		return "", err
	}

	ok := utils.CheckPassword(user.PasswordHash, password, s.cfg.PasswordSalt, user.PasswordSalt)
	if !ok {
		return "", storage.ErrUnauthorized
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.cfg.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	return t.SignedString([]byte(s.cfg.TokenSigningKey))
}

func (s *Service) ParseToken(ctx context.Context, bearer string) (user storage.User, err error) {
	token, err := jwt.ParseWithClaims(bearer, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(s.cfg.TokenSigningKey), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return storage.User{}, ErrWrongTokenClaimsType
	}

	user, err = s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return storage.User{}, ErrUnauthorized
	}

	err = s.UpdateUserLastOnline(ctx, user)
	if err != nil {
		return storage.User{}, ErrUpdateUserLastOnline
	}

	return
}

func (s *Service) AuthDevice(ctx context.Context, deviceID uuid.UUID, token string) (device storage.Device, err error) {
	device, err = s.CheckDeviceAuthorization(ctx, deviceID, token)
	if err != nil {
		return storage.Device{}, err
	}

	err = s.UpdateDeviceLastOnline(ctx, device)
	if err != nil {
		return storage.Device{}, err
	}

	return device, nil
}
