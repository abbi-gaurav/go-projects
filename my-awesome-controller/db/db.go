package db

import (
	"github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"strings"
)

type DB interface {
	Add(fqn string, cake *v1.Cake)
	Update(fqn string, cake *v1.Cake)
	Delete(fqn string)
	Get(fqn string) *v1.Cake
}

func New(dbType string) DB {
	switch strings.ToLower(dbType) {
	case "memory":
		return newInMemory()
	default:
		panic("Only in memory db implemented")
	}
}
