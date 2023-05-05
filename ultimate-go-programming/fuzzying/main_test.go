package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func FuzzHandler(f *testing.F) {
	tests := []struct {
		input  string
		result string
	}{
		{input: `{"x": 4, "y": 2}`, result: `{"result":2}`},
		{input: `{"x": 6, "y": 2}`, result: `{"result":3}`},
	}
	for _, tt := range tests {
		f.Add(tt.input)
	}

	f.Fuzz(func(t *testing.T, input string) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080",
			strings.NewReader(input))
		handler(w, req)
	})
}
