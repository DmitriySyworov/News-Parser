package handler_response

import (
	"encoding/json"
	"net/http"
)

func HandlerResponse(writer http.ResponseWriter, v any, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	errEncode := json.NewEncoder(writer).Encode(v)
	if errEncode != nil {
		panic(errEncode)
	}
}
