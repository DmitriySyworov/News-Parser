package user

import "time"

type RequestUpdateUser struct {
	Name        string `validate:"omitempty,min=2,max=64"`
	NewEmail    string `validate:"omitempty,email,required_with=Password"`
	Password    string `validate:"omitempty,min=8,max=24"`
	NewPassword string `json:"new-password" validate:"omitempty,required_with=Password,min=8,max=24"`
}
type RequestRemoveOrDelete struct {
	Password string `validate:"required,min=8,max=24"`
}
type ResponseUser struct {
	CreatedAt time.Time
	Name      string
	Email     string
	UserUUID  string
}
