package xerror

import (
	"fmt"
)

type StatusError struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

func (e *StatusError) WithDetails(details ...interface{}) {
	e.Details = append(e.Details, details...)
}

func (e *StatusError) Is(target error) bool {
	err, ok := target.(*StatusError)
	if ok {
		return e.Code == err.Code
	}
	return false
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", e.Code, e.Message, e.Details)
}
