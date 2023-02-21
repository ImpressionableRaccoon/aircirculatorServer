package storage

import (
	"errors"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyNoPermissions = errors.New("company exists, but you have not permission to access it")
	ErrDeviceAlreadyExists  = errors.New("device already exists")
	ErrDeviceNotFound       = errors.New("device not found")
)
