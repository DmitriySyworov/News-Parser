package user

import "errors"

var (
	ErrIncorrectType    = errors.New("the 'type' must be a hard-delete or soft-delete")
	ErrUpdateUser       = errors.New("failed to update user")
	ErrIncorrectAction  = errors.New("the 'action' must be a soft-delete, hard-delete or update")
	ErrFailedDeleteUser = errors.New("failed to delete user")
	ErrFailedRemoveUser = errors.New("failed to remove user")
)
