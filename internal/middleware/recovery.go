package middleware

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/response"
	"net/http"
)

func (m *ManagerMiddleware) RecoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				defer func() {
					m.resp = response.Response{}
					m.ContextValues = &ContextValues{
						DataLog: loggers.DataLog{
							MapLog: make(map[string]any),
						},
					}
				}()
				m.resp.Errors = append(m.resp.Errors, response.Error{
					Message: custom_errors.ErrCriticalServer.Error(),
					Status:  http.StatusInternalServerError,
				})
				m.Logger.HandlerLogger(request.Pattern, m.DataLog.UserUUID, 500, m.DataLog.Errors, m.DataLog.MapLog)
				response.HandlerResponse(writer, m.resp, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
