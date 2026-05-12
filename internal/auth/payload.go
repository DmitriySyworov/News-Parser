package auth

type RequestRegister struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
type RequestLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type ResponseConfirm struct {
	JWT string
}
