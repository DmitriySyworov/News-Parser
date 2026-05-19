package auth

import "errors"

var (
	ErrSaveDataUser            = errors.New("failed to save data user")
	ErrLoginEmailOrPassword    = errors.New("incorrect email or password")
	ErrIncorrectActionRecovery = errors.New("the 'action' must be a recovery-password or recovery-remove")
	ErrIncorrectActionConfirm  = errors.New("the 'action' must be a register, recovery-password or recovery-remove")
	ErrNotNewPassword          = errors.New("when recovering your password, you must specify a new one")
	ErrFailedChangePassword    = errors.New("failed to change password")
	ErrFailedRecovery          = errors.New("failed to recovery")
)
