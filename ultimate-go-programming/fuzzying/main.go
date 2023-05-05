package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

type request struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type response struct {
	Result int `json:"result"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var req *request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "unable to decode data", http.StatusInternalServerError)
		return
	}
	if req.Y == 0 {
		http.Error(w, "unable to divide by zero", http.StatusBadRequest)
		return
	}

	var resp response
	resp.Result = req.X / req.Y
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "unable to decode data", http.StatusInternalServerError)
		return
	}
}
