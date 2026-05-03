package custom_errors

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrIncorrectDate  = errors.New("incorrect format date or empty")
)
