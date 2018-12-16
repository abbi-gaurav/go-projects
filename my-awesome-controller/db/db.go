package db

import "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"

type DB interface {
	Add(fqn string, cake *v1.Cake)
	Update(fqn string, cake *v1.Cake)
	Delete(fqn string)
	Get(fqn string) *v1.Cake
}

func New(dbType string) DB {
	switch dbType {
	case "memory":
		return newInMemory()
	default:
		panic("Only in memory db implemented")
	}
}
