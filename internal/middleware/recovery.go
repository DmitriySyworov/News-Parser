package middleware

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/response"
	"log"
	"net/http"
)

func (m *ManagerMiddleware) RecoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				log.Println(errPanic)
				m.resp.Errors = append(m.resp.Errors, response.Error{
					Message: custom_errors.ErrCriticalServer.Error(),
					Status:  http.StatusInternalServerError,
				})
				response.HandlerResponse(writer, m.resp, http.StatusInternalServerError)
				m.resp = response.Response[any]{}
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
