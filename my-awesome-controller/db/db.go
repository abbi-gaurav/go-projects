package db

import "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"

type DB interface {
	Add(cake *v1.Cake)
	Update(cake *v1.Cake)
	Delete(cake *v1.Cake)
}

func fqName(cake *v1.Cake) string {
	return cake.Namespace + "." + cake.Name
}

func New(dbType string) DB {
	switch dbType {
	case "memory":
		return newInMemory()
	default:
		panic("Only in memory db implemented")
	}
}
