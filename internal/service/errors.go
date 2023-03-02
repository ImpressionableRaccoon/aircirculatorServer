package service

import (
	"errors"
)

var (
	ErrUnauthorized             = errors.New("unauthorized")
	ErrUpdateUserLastOnline     = errors.New("error update user last online")
	ErrInvalidSigningMethod     = errors.New("invalid signing method")
	ErrLoginIsEmpty             = errors.New("login is empty")
	ErrPasswordIsEmpty          = errors.New("password is empty")
	ErrWrongTokenClaimsType     = errors.New("wrong token claims type")
	ErrWrongDeviceID            = errors.New("wrong device id")
	ErrWrongDeviceFormat        = errors.New("wrong device format")
	ErrWrongScheduleID          = errors.New("wrong schedule id")
	ErrWrongScheduleWeek        = errors.New("wrong schedule week day")
	ErrWrongScheduleTimeStart   = errors.New("wrong schedule time start")
	ErrWrongScheduleTimeStop    = errors.New("wrong schedule time stop")
	ErrTimeStopNotMoreTimeStart = errors.New("wrong schedule: time stop must be later than time start")
	ErrFirmwareNoUpdates        = errors.New("no firmware updates")
)
