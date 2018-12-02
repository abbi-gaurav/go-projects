package router

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"os"
	"time"
	"github.com/go-learning-projects/k8s-ready-service/model"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	actualInfo := model.Info{
		BuildTime: time.Now().Format("20060102_03:04:05"),
		Commit:    "some test hash",
		Release:   "0.0.8",
	}

	router := Router(actualInfo.BuildTime, actualInfo.Commit, actualInfo.Release)
	ts = httptest.NewServer(router)

	ret := m.Run()
	ts.Close()

	os.Exit(ret)
}

func verifyStatusCode(res *http.Response, expectedStatusCode int, t *testing.T) {
	if res.StatusCode != expectedStatusCode {
		t.Errorf("Status code is wrong, have: %d, want: %d", res.StatusCode, expectedStatusCode)
	}
}

func checkIfError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestRouterCorrectPathAndMethod(t *testing.T) {
	res, err := http.Get(ts.URL + "/home")
	checkIfError(err, t)
	verifyStatusCode(res, http.StatusOK, t)
}

func TestRouterCorrectPathWrongMethod(t *testing.T) {
	res, err := http.Post(ts.URL+"/home", "text/plain", nil)
	checkIfError(err, t)

	verifyStatusCode(res, http.StatusMethodNotAllowed, t)
}

func TestRouterWrongPath(t *testing.T) {
	res, err := http.Get(ts.URL + "/not-exists")
	checkIfError(err, t)

	verifyStatusCode(res, http.StatusNotFound, t)
}
