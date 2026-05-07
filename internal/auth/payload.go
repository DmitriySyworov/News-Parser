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
type ResponseAuth struct {
	Message string `json:"message"`
	JWTTemp string `json:"jwt-temp"`
	Error   string `json:"error"`
}
type ResponseConfirm struct {
	JWT   string
	Error string
}
type RequestConfirm struct {
	Code uint `json:"code" validate:"required"`
}
