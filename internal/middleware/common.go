package middleware

import (
	"app/news-parser/internal/custom_errors"
)

type ManagerMiddleware struct {
	Signature string
	respError custom_errors.ResponseError
	ContextToken
}
type ContextToken struct {
	SessionID string
	UUID      string
}

func NewManagerMiddleware(signature string) *ManagerMiddleware {
	return &ManagerMiddleware{
		Signature: signature,
	}
}
