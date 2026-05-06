package auth

import "errors"

var (
	ErrSendLetter           = errors.New("we were unable to send a letter to the specified email")
	ErrFailedSecurity       = errors.New("failed to secure the authorization session")
	ErrExpiredSession       = errors.New("session time has expired")
	ErrIncorrectCode        = errors.New("the code is incorrect")
	ErrSaveDataUser         = errors.New("failed to save data user")
	ErrLoginEmailOrPassword = errors.New("incorrect email or password")
)
