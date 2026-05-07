package custom_errors

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrIncorrectDate   = errors.New("incorrect format date or empty")
	ErrUserExist       = errors.New("this user already exists")
	ErrIncorrectToken  = errors.New("incorrect token")
	ErrIncorrectAction = errors.New("no such action exists")
ErrUserNotFound = errors.New("no such user record found")
)
