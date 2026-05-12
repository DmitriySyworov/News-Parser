package auth

import "errors"

var (
	ErrSaveDataUser         = errors.New("failed to save data user")
	ErrLoginEmailOrPassword = errors.New("incorrect email or password")
)
