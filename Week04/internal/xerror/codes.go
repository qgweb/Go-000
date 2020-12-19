package xerror

import (
	"errors"
	"fmt"
)

// HTTP Mapping: 401 Unauthorized
func Unauthorized(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    16,
		Message: "Unauthenticated",
		Details: []interface{}{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

func IsUnauthorized(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 16
	}
	return false
}

// HTTP Mapping: 404 Not Found
func NotFound(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    5,
		Message: "NotFound",
		Details: []interface{}{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

func IsNotFound(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 5
	}
	return false
}

// HTTP Mapping: 500 Internal Server Error
func Internal(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    13,
		Message: "Internal",
		Details: []interface{}{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

func IsInternal(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 13
	}
	return false
}
