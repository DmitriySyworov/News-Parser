package middleware

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/JWT"
	"app/news-parser/pkg/handler_response"
	"context"
	"net/http"
	"strings"
)

func helperValidateToken(header string) (string, error) {
	sliceHeader := strings.Split(header, " ")
	if len(sliceHeader) != 2 {
		return "", custom_errors.ErrIncorrectToken
	}
	if strings.Count(sliceHeader[1], ".") != 2 {
		return "", custom_errors.ErrIncorrectToken
	}
	return sliceHeader[1], nil
}

const (
	KeySessionJWT = "keySessionJWT"
	KeyAuthToken  = "keyAuthToken"
)

func (m *ManagerMiddleware)IsTemporaryJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.resp.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, m.resp.Error, http.StatusBadRequest)
			return
		}
		j := JWT.NewJWT(m.Signature)
		sessionId, errParseJwt := j.ParseTemporaryJWT(token)
		if errParseJwt != nil {
			m.resp.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, m.resp.Error, http.StatusBadRequest)
			return
		}
		valueCtx := context.WithValue(context.Background(), KeySessionJWT, sessionId)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
func (m *ManagerMiddleware)IsAuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.resp.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, m.resp.Error, http.StatusBadRequest)
			return
		}
		j := JWT.NewJWT(m.Signature)
		UUID, errParseJwt := j.ParseJWT(token)
		if errParseJwt != nil {
			m.resp.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, m.resp.Error, http.StatusBadRequest)
			return
		}
		valueCtx := context.WithValue(context.Background(), KeyAuthToken, UUID)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
