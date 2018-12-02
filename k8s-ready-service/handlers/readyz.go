package handlers

import (
	"sync/atomic"
	"net/http"
)

func Readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(responseWriter, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
}
