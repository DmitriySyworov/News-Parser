package custom_errors

import "errors"

type ResponseError struct {
	Success bool
	Errors  []Error
}
type Error struct {
	Message string
	Status  int
}

var (
	ErrIncorrectCode        = errors.New("the code is incorrect")
	ErrSession            = errors.New("session is empty or time has expired")
	ErrFailedSecurity     = errors.New("failed to secure the authorization session")
	ErrSendLetter         = errors.New("we were unable to send a letter to the specified email")
	ErrIncorrectArticleId = errors.New("incorrect article_ID format or empty")
	ErrIncorrectWithText  = errors.New("the 'withText' must be a boolean is true or false")
	ErrRecordNotFound     = errors.New("record not found")
	ErrIncorrectDate      = errors.New("incorrect format date or empty")
	ErrUserExist          = errors.New("this user already exists")
	ErrIncorrectToken     = errors.New("incorrect token")
	ErrIncorrectAction    = errors.New("no such action exists")
	ErrUserNotExist       = errors.New("such user does not exist")
	ErrIncorrectOffset    = errors.New("the 'offset' must be a positive integer")
	ErrIncorrectLimit     = errors.New("the 'limit' must be a positive integer")
)
