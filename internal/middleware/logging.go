package middleware

import (
	"context"
	"net/http"
)

func (m *ManagerMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		wrapperWriter := &WrapperWriter{
			StatusCode:     http.StatusOK,
			ResponseWriter: writer,
		}
		ctxValue := context.WithValue(context.Background(), KeyContextValues, m.ContextValues)
		ctxRequest := request.WithContext(ctxValue)
		next.ServeHTTP(wrapperWriter, ctxRequest)
		m.Logger.HandlerLogger(ctxRequest.Pattern, m.Logger.DataLog.UserUUID, wrapperWriter.StatusCode, m.Logger.Errors, m.Logger.DataLog.MapLog)
	})
}

type WrapperWriter struct {
	StatusCode int
	http.ResponseWriter
}

func (w *WrapperWriter) WriteHeader(status int) {
	w.StatusCode = status
	w.ResponseWriter.WriteHeader(status)
}
