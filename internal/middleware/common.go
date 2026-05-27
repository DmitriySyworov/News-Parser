package middleware

import (
	"app/news-parser/internal/response"
)

type ManagerMiddleware struct {
	Signature string
	resp      response.Response[any]
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
