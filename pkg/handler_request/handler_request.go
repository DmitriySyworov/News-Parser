package handler_request

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	ErrIncorrectFormat = errors.New("incorrect format of transmitted data")
	ErrInvalidData     = errors.New("the data transmitted was incorrect")
)

func HandlerRequest[T any](request *http.Request) (*T, error) {
	var payload T
	errDecode := json.NewDecoder(request.Body).Decode(&payload)
	if errDecode != nil {
		return nil, ErrIncorrectFormat
	}
	errValidate := validator.New().Struct(payload)
	if errValidate != nil {
		return nil, ErrInvalidData
	}
	return &payload, nil
}
