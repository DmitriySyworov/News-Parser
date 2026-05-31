package middleware

import (
	"app/news-parser/internal/JWT"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/response"
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

func (m *ManagerMiddleware) IsTemporaryJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			m.resp = response.Response{}
		}()
		header := request.Header.Get("X-Temp-Token")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.resp.Errors = append(m.resp.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, m.resp, http.StatusUnauthorized)
			return
		}
		j := JWT.NewJWT(m.Signature)
		sessionId, errParseJwt := j.ParseTemporaryJWT(token)
		if errParseJwt != nil {
			m.resp.Errors = append(m.resp.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, m.resp, http.StatusUnauthorized)
			return
		}
		if tokens, ok := request.Context().Value(KeyContextValues).(ContextValues); ok && tokens.UserUUID != "" {
			m.ContextValues.UserUUID = tokens.UserUUID
		}
		m.ContextValues.SessionID = sessionId
		valueCtx := context.WithValue(context.Background(), KeyContextValues, m.ContextValues)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
func (m *ManagerMiddleware) IsAuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			m.resp = response.Response{}
		}()
		header := request.Header.Get("Authorization")
		token, errToken := helperValidateToken(header)
		if errToken != nil {
			m.resp.Errors = append(m.resp.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, m.resp, http.StatusUnauthorized)
			return
		}
		j := JWT.NewJWT(m.Signature)
		UUID, errParseJwt := j.ParseJWT(token)

		if errParseJwt != nil {
			m.resp.Errors = append(m.resp.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, m.resp, http.StatusUnauthorized)
			return
		}
		if tokens, ok := request.Context().Value(KeyContextValues).(ContextValues); ok && tokens.SessionID != "" {
			m.ContextValues.SessionID = tokens.SessionID
		}
		m.ContextValues.UserUUID = UUID
		valueCtx := context.WithValue(context.Background(), KeyContextValues, m.ContextValues)
		requestCTX := request.WithContext(valueCtx)
		next.ServeHTTP(writer, requestCTX)
	})
}
