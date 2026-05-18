package middleware

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/handler_response"
	"log"
	"net/http"
)

func (m *ManagerMiddleware) RecoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				log.Println(errPanic)
				m.respError.Errors = append(m.respError.Errors, custom_errors.Error{
					Message: custom_errors.ErrCriticalServer.Error(),
					Status:  http.StatusInternalServerError,
				})
				handler_response.HandlerResponse(writer, m.respError, http.StatusInternalServerError)
				m.respError = custom_errors.ResponseError{}
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
