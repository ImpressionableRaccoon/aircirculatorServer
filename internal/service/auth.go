package service

import (
	"context"
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
	hash, salt, err := utils.PreparePassword(password, s.cfg.PasswordSalt)
	if err != nil {
		return "", err
	}

	err = s.st.SignUp(ctx, login, hash, salt)
	if err != nil {
		return "", err
	}

	return s.SignIn(ctx, login, password)
}

func (s *Service) SignIn(ctx context.Context, login string, password string) (token string, err error) {
	user, err := s.st.GetUserByLogin(ctx, login)

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

	return t.SignedString(s.cfg.TokenSigningKey)
}

func (s *Service) ParseToken(ctx context.Context, bearer string) (user storage.User, err error) {
	token, err := jwt.ParseWithClaims(bearer, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return s.cfg.TokenSigningKey, nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return storage.User{}, ErrWrongTokenClaimsType
	}

	user, err = s.st.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return storage.User{}, ErrUnauthorized
	}

	return
}
