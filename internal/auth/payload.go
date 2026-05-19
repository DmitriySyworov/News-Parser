package auth

type RequestRegister struct {
	Name     string `validate:"required,min=2,max=64"`
	Email    string `validate:"required,email,min=5,max=256"`
	Password string `validate:"required,min=8,max=24"`
}
type RequestLogin struct {
	Email    string `validate:"required,email,min=5,max=256"`
	Password string `validate:"required"`
}

type ResponseConfirm struct {
	JWT string
}

type RequestRecovery struct {
	Email       string `validate:"required,min=5,max=256"`
	NewPassword string `json:"new-password" validate:"omitempty,min=8,max=24"`
}
