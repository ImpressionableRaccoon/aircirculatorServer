package storage

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
)
