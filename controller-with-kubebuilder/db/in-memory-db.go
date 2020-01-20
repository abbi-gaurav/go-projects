package db

import (
	"sync"

	"github.com/abbi-gaurav/go-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
)

type inMemoryDB struct {
	sm sync.Map
}

func (m *inMemoryDB) Add(fqn string, sloop *v1beta1.SloopSpec) {
	m.sm.Store(fqn, sloop)
}

func (m *inMemoryDB) Update(fqn string, sloop *v1beta1.SloopSpec) {
	m.sm.Store(fqn, sloop)
}

func (m *inMemoryDB) Delete(fqn string) {
	m.sm.Delete(fqn)
}

func (m *inMemoryDB) Get(fqn string) *v1beta1.SloopSpec {
	obj, _ := m.sm.Load(fqn)

	if obj == nil {
		return nil
	}

	return obj.(*v1beta1.SloopSpec)
}

func newInMemory() DB {
	return &inMemoryDB{}
}
