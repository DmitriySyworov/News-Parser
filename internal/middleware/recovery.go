package middleware

import (
	"app/news-parser/pkg/handler_response"
	"errors"
	"log"
	"net/http"
)

func (m *ManagerMiddleware) RecoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				log.Println(errPanic)
				m.resp.Error = errors.New("critical error on the server side").Error()
				handler_response.HandlerResponse(writer, m.resp, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
