package handlers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/go-learning-projects/k8s-ready-service/model"
	"time"
)

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()

	actualInfo := model.Info{
		BuildTime: time.Now().Format("20060102_03:04:05"),
		Commit:    "some test hash",
		Release:   "0.0.8",
	}

	handler := Home(actualInfo.BuildTime, actualInfo.Commit, actualInfo.Release)
	handler(w, nil)

	resp := w.Result()

	if have, want := resp.StatusCode, http.StatusOK; have != want {
		t.Errorf("Status code is wrong have:%d, want:%d", have, want)
	}

	greeting, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	info := model.Info{}

	if err = json.Unmarshal(greeting, &info); err != nil {
		t.Fatal(err)
	}

	if have, want := actualInfo, info; have != want {
		t.Errorf("The response is wrong, have:%v, want:%v", have, want)
	}
}
