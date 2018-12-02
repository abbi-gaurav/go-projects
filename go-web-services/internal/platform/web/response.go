package web

import (
	"github.com/pkg/errors"
	"context"
	"net/http"
	"encoding/json"
	"log"
	"io"
)

type JSONError struct {
	Error string `json:"error"`
}

var (
	ErrNotFound          = errors.New("Entity Not Found")
	ErrorDBNotConfigured = errors.New("DB Not configured")
)

func Error(ctx context.Context, w http.ResponseWriter, err error) {
	switch errors.Cause(err) {
	case ErrNotFound:
		RespondError(ctx, w, err, http.StatusNotFound)
		return
	}

	RespondError(ctx, w, err, http.StatusInternalServerError)
}
func RespondError(ctx context.Context, responseWriter http.ResponseWriter, err error, statusCode int) {
	Respond(ctx, responseWriter, JSONError{Error: err.Error()}, statusCode)
}
func Respond(ctx context.Context, writer http.ResponseWriter, data interface{}, statusCode int) {
	v := ctx.Value(KeyValues).(*Values)
	v.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		writer.WriteHeader(statusCode)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Eror marshalling %v", err)
		jsonData = []byte("{}")
	}

	io.WriteString(writer, string(jsonData))
}
