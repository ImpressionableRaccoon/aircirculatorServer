package storage

import (
	"errors"
)

var (
	ErrUserAlreadyExists error = errors.New("user already exists")
	ErrUnauthorized      error = errors.New("unauthorized")
)
