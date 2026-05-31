package middleware

import (
	"app/news-parser/internal/loggers"
	"context"
	"net/http"
)

func (m *ManagerMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			m.ContextValues = &ContextValues{
				DataLog: loggers.DataLog{
					MapLog: make(map[string]any),
				},
			}
		}()
		wrapperWriter := &WrapperWriter{
			StatusCode:     http.StatusOK,
			ResponseWriter: writer,
		}
		ctxValue := context.WithValue(context.Background(), KeyContextValues, m.ContextValues)
		ctxRequest := request.WithContext(ctxValue)
		next.ServeHTTP(wrapperWriter, ctxRequest)
		m.Logger.HandlerLogger(ctxRequest.Pattern, m.DataLog.UserUUID, wrapperWriter.StatusCode, m.DataLog.Errors, m.DataLog.MapLog)
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
