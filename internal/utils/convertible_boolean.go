package utils

import (
	"errors"
)

var (
	ErrWrongBooleanFormat = errors.New("boolean unmarshal error")
)

type ConvertibleBoolean bool

func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		*bit = true
	} else if asString == "0" || asString == "false" {
		*bit = false
	} else {
		return ErrWrongBooleanFormat
	}
	return nil
}
