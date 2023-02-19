package service

import (
	"errors"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrWrongTokenClaimsType = errors.New("wrong token claims type")
)
