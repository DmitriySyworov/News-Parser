package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool
	Data    any
	Errors  []Error
}
type Error struct {
	Message string
	Status  int
}

func HandlerResponse(writer http.ResponseWriter, resp Response, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	errEncode := json.NewEncoder(writer).Encode(resp)
	if errEncode != nil {
		panic(errEncode)
	}
}
