package service

import (
	"errors"
)

var (
	ErrUnauthorized             = errors.New("unauthorized")
	ErrInvalidSigningMethod     = errors.New("invalid signing method")
	ErrWrongTokenClaimsType     = errors.New("wrong token claims type")
	ErrWrongDeviceFormat        = errors.New("wrong device format")
	ErrWrongScheduleID          = errors.New("wrong schedule id")
	ErrWrongScheduleWeek        = errors.New("wrong schedule week day")
	ErrWrongScheduleTimeStart   = errors.New("wrong schedule time start")
	ErrWrongScheduleTimeStop    = errors.New("wrong schedule time stop")
	ErrTimeStopNotMoreTimeStart = errors.New("wrong schedule: time stop must be later than time start")
)
