package handler_request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	ErrIncorrectFormat = errors.New("incorrect format of transmitted data")
	ErrInvalidData     = errors.New("the data transmitted was incorrect")
	ErrBodyIsEmpty     = errors.New("request body is empty")
)

func HandlerRequest[T any](request *http.Request) (*T, error) {
	var payload T
	data, errRead := io.ReadAll(request.Body)
	if errRead != nil || len(data) < 3 {
		return nil, ErrBodyIsEmpty
	}
	errDecode := json.Unmarshal(data, &payload)
	if errDecode != nil {
		return nil, ErrIncorrectFormat
	}
	errValidate := validator.New().Struct(payload)
	if errValidate != nil {
		return nil, ErrInvalidData
	}
	return &payload, nil
}
