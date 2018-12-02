package handlers

import (
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/mid"
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/platform/db"
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/platform/web"
	"net/http"
)

func API(masterDB *db.DB) http.Handler {
	app := web.New(mid.RequestLogger, mid.ErrorHandler)

	u := User{
		MasterDB: masterDB,
	}
	app.Handle("GET", "/v1/users", u.List)

	return app
}
