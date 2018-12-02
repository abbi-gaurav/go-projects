package handlers

import (
	"encoding/json"
	"github.com/abbi-gaurav/go-learning-projects/k8s-ready-service/model"
	"log"
	"net/http"
)

func Home(buildTime, commit, release string) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		info := model.Info{
			BuildTime: buildTime,
			Commit:    commit,
			Release:   release,
		}

		body, err := json.Marshal(info)
		if err != nil {
			log.Printf("Could not encode info data %v", err)
			http.Error(responseWriter, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write(body)
	}
}
