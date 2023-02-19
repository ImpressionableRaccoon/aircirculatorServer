package service

import (
	"errors"
)

var (
	ErrUnauthorized         error = errors.New("unauthorized")
	ErrInvalidSigningMethod error = errors.New("invalid signing method")
	ErrWrongTokenClaimsType error = errors.New("wrong token claims type")
)
