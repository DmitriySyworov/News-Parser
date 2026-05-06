package JWT

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	JwtSignature []byte
}

func NewJWT(signature string) *Jwt {
	return &Jwt{
		JwtSignature: []byte(signature),
	}
}
func (j *Jwt) CreateTemporaryJWT(sessionId string) (string, error) {
	jCreate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_session": sessionId,
		"expires_at": time.Now().Add(5 * time.Minute),
	})
	token, errToken := jCreate.SignedString(j.JwtSignature)
	if errToken != nil {
		return "", errToken
	}
	return token, nil
}

func (j *Jwt) ParseTemporaryJWT(token string) (string, error) {
	valueToken, errParse := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.JwtSignature, nil
	})
	if errParse != nil {
		return "", errParse
	}
	session, ok := valueToken.Claims.(jwt.MapClaims)["id_session"].(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}
	return session, nil
}
func (j *Jwt) CreateJWT(UUID string) (string, error) {
	jCreate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UUID": UUID,
	})
	token, errToken := jCreate.SignedString(j.JwtSignature)
	if errToken != nil {
		return "", errToken
	}
	return token, nil
}

func (j *Jwt) ParseJWT(token string) (string, error) {
	valueToken, errParse := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.JwtSignature, nil
	})
	if errParse != nil {
		return "", errParse
	}
	UUID, ok := valueToken.Claims.(jwt.MapClaims)["UUID"].(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}
	return UUID, nil
}
