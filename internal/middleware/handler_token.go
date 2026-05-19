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
	KeyContext = "keyContext"
)

func (m *ManagerMiddleware) IsTemporaryJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			m.respError = custom_errors.ResponseError{}
		}()
		header := request.Header.Get("X-Temp-Token")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.respError.Errors = append(m.respError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, m.respError, http.StatusUnauthorized)
			return
		}
		j := JWT.NewJWT(m.Signature)
		sessionId, errParseJwt := j.ParseTemporaryJWT(token)
		if errParseJwt != nil {
			m.respError.Errors = append(m.respError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, m.respError, http.StatusUnauthorized)
			return
		}
		if tokens, ok := request.Context().Value(KeyContext).(ContextToken); ok && tokens.UUID != "" {
			m.ContextToken.UUID = tokens.UUID
		}
		m.ContextToken.SessionID = sessionId
		valueCtx := context.WithValue(context.Background(), KeyContext, m.ContextToken)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
func (m *ManagerMiddleware) IsAuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			m.respError = custom_errors.ResponseError{}
		}()
		header := request.Header.Get("Authorization")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.respError.Errors = append(m.respError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, m.respError, http.StatusUnauthorized)
			return
		}
		j := JWT.NewJWT(m.Signature)
		UUID, errParseJwt := j.ParseJWT(token)

		if errParseJwt != nil {
			m.respError.Errors = append(m.respError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, m.respError, http.StatusUnauthorized)
			return
		}
		if tokens, ok := request.Context().Value(KeyContext).(ContextToken); ok && tokens.SessionID != "" {
			m.ContextToken.SessionID = tokens.SessionID
		}
		m.ContextToken.UUID = UUID
		valueCtx := context.WithValue(context.Background(), KeyContext, m.ContextToken)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
