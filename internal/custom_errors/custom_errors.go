package custom_errors

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrIncorrectDate           = errors.New("incorrect format date or empty")
	ErrUserExist               = errors.New("this user already exists")
	ErrIncorrectToken          = errors.New("incorrect token")
	ErrIncorrectAction         = errors.New("no such action exists")
	ErrUserNotExist            = errors.New("such user does not exist")
	ErrIncorrectOffset         = errors.New("the 'offset' must be a positive integer")
	ErrIncorrectLimit          = errors.New("the 'limit' must be a positive integer")
	ErrIncorrectOffsetAndLimit = errors.New("the 'offset' must be a positive integer. The 'limit' must be a positive integer")
)
