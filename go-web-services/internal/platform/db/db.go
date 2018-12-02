package db

import (
	"gopkg.in/mgo.v2"
	"github.com/pkg/errors"
	"time"
	"context"
	"encoding/json"
)

var ErrInvalidDBProvided = errors.New("Invalid DB provided")

type DB struct {
	database *mgo.Database
	session  *mgo.Session
}

func New(url string, timeout time.Duration) (*DB, error) {

	if timeout == 0 {
		timeout = 60 * time.Second
	}

	sess, err := mgo.DialWithTimeout(url, timeout)

	if err != nil {
		return nil, errors.Wrapf(err, "mgo.DialWithTimeout : %s, %v", url, timeout)
	}

	sess.SetMode(mgo.Monotonic, true)

	db := DB{
		database: sess.DB(""),
		session:  sess,
	}

	return &db, nil
}

func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) Copy() (*DB, error) {
	sess := db.session.Copy()

	newDB := DB{
		database: sess.DB(""),
		session:  sess,
	}

	return &newDB, nil
}

func (db *DB) Execute(ctx context.Context, collName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil || db.session == nil")
	}

	return f(db.database.C(collName))
}

func JSONString(value interface{}) string {
	js, err := json.Marshal(value)

	if err != nil {
		return ""
	}

	return string(js)
}
