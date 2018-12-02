package handlers_test

import (
	"github.com/go-learning-projects/ultimate-go-programming/testing/endpoint_test/handlers"
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func init() {
	handlers.Routes()
}

func TestSendJSON(t *testing.T) {
	t.Log("Given the need to test /sendjson ep")
	{
		req, err := http.NewRequest("GET", "/sendjson", nil)
		if err != nil {
			t.Fatal("\tShould be able to create a request.", ballotX, err)
		}

		t.Log("\t Should be able to create a request", checkMark)

		rw := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rw, req)

		if rw.Code != 200 {
			t.Fatal("\tShould receive 200", ballotX, rw.Code)
		}
		t.Log("\tShould receive 200", checkMark)

		u := struct {
			Name  string
			Email string
		}{}

		if err = json.NewDecoder(rw.Body).Decode(&u); err != nil {
			t.Fatal("\tShould decomopose the response", ballotX)
		}
		t.Log("\tShould decomopose response", checkMark)

		if u.Name == "ss" {
			t.Log("\tShould have correct name", checkMark)
		} else {
			t.Error("\tShould have name", ballotX, u.Name)
		}
	}
}
