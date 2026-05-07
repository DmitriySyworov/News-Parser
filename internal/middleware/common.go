package middleware

type ManagerMiddleware struct{
	Signature string
	 resp *responseMiddlewareError
}
type responseMiddlewareError struct {
	Error string
}
func NewManagerMiddleware(signature string)*ManagerMiddleware{
	return &ManagerMiddleware{
		Signature: signature,
	}
}