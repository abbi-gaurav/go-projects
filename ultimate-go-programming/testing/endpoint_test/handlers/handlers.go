package handlers

import (
	"net/http"
	"encoding/json"
)

func Routes() {
	http.HandleFunc("/sendjson", sendJSON)
}

func sendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name  string
		Email string
	}{
		Name:  "ss",
		Email: "s@s.com",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(&u)

}
