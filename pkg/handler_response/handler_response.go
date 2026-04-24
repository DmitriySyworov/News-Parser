package handler_response

import (
	"encoding/json"
	"net/http"
)

func HandlerResponse(writer http.ResponseWriter, v any, status int) {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(writer).Encode(v)
	if errEncode != nil {
		panic(errEncode)
	}
}
