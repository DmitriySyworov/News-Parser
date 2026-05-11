package middleware

import "app/news-parser/internal/custom_errors"

type ManagerMiddleware struct {
	Signature string
	respError custom_errors.ResponseError
}

func NewManagerMiddleware(signature string) *ManagerMiddleware {
	return &ManagerMiddleware{
		Signature: signature,
	}
}
