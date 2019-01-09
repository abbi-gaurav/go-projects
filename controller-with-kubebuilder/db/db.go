package db

import (
	"strings"

	"github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
)

type DB interface {
	Add(fqn string, sloop *v1beta1.SloopSpec)
	Update(fqn string, sloop *v1beta1.SloopSpec)
	Delete(fqn string)
	Get(fqn string) *v1beta1.SloopSpec
}

func New(dbType string) DB {
	switch strings.ToLower(dbType) {
	case "memory":
		return newInMemory()
	default:
		panic("no other than in memory database")
	}
}
