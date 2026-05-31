package custom_errors

import "errors"

var (
	ErrIncorrectCode           = errors.New("the code is incorrect")
	ErrSession                 = errors.New("session is empty or time has expired")
	ErrFailedSecurity          = errors.New("failed to secure the authorization session")
	ErrSendLetter              = errors.New("we were unable to send a letter to the specified email")
	ErrIncorrectWithText       = errors.New("the 'withText' must be a boolean is true or false")
	ErrRecordNotFound          = errors.New("record not found")
	ErrIncorrectDate           = errors.New("the date format must be 'YYYY-MM-DD'")
	ErrUserExist               = errors.New("this user already exists")
	ErrIncorrectToken          = errors.New("incorrect token")
	ErrUserNotExist            = errors.New("such user does not exist")
	ErrIncorrectOffset         = errors.New("the 'offset' must be a positive integer")
	ErrIncorrectLimit          = errors.New("the 'limit' must be a positive integer")
	ErrCriticalServer          = errors.New("critical error on the server side")
	ErrFailedTypeContextValues = errors.New("failed type assertion ContextValues: ")
	ErrIncorrectPassword       = errors.New("the password must be between 8 and 24 characters long")
	ErrIncorrectNewPassword    = errors.New("the new_password must be between 8 and 24 characters long")
	ErrIncorrectName           = errors.New("the name must be between 2 and 64 characters long")
	ErrIncorrectEmail          = errors.New("incorrect email")
)
