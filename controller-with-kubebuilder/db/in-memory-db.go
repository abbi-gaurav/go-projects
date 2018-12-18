package db

import (
	"github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
	"sync"
)

type inMemoryDB struct {
	sm sync.Map
}

func (m *inMemoryDB) Add(fqn string, sloop *v1beta1.Sloop) {
	m.sm.Store(fqn, sloop)
}

func (m *inMemoryDB) Update(fqn string, sloop *v1beta1.Sloop) {
	m.sm.Store(fqn, sloop)
}

func (m *inMemoryDB) Delete(fqn string) {
	m.sm.Delete(fqn)
}

func (m *inMemoryDB) Get(fqn string) *v1beta1.Sloop {
	obj, _ := m.sm.Load(fqn)

	if obj == nil {
		return nil
	}

	return obj.(*v1beta1.Sloop)
}

func newInMemory() DB {
	return &inMemoryDB{}
}
