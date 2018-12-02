package handlers

import (
	"context"
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/platform/db"
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/platform/web"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"net/http"
)

type User struct {
	MasterDB *db.DB
}

func (u *User) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	reqDB, err := u.MasterDB.Copy()
	if err != nil {
		return errors.Wrapf(web.ErrorDBNotConfigured, "")
	}
	defer reqDB.Close()

	data := struct {
		Name  string
		Email string
	}{
		Name:  "Gaurav",
		Email: "abc@xyz.com",
	}

	f := func(collection *mgo.Collection) error {
		return collection.Insert(data)
	}

	if err := reqDB.Execute(ctx, "users", f); err != nil {
		return errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.JSONString(data)))
	}

	web.Respond(ctx, w, data, http.StatusOK)
	return nil
}
