package user

import "time"

type RequestUpdateUser struct {
	Name        string
	NewEmail    string `validate:"omitempty,email,required_with=Password"`
	Password    string
	NewPassword string `validate:"omitempty,required_with=Password"`
}
type RequestRemoveOrDelete struct {
	Password string `validate:"required"`
}
type ResponseUser struct {
	CreatedAt time.Time
	Name      string
	Email     string
	UUIDUser  string
}
