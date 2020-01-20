package db

import (
	"github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"sync"
)

type inMemoryDB struct {
	sm sync.Map
}

func newInMemory() DB {
	return &inMemoryDB{}
}

func (memDB *inMemoryDB) Add(fqn string, cake *v1.Cake) {
	memDB.set(fqn, cake, "add")
}

func (memDB *inMemoryDB) Update(fqn string, cake *v1.Cake) {
	memDB.set(fqn, cake, "update")
}

func (memDB *inMemoryDB) Delete(fqn string) {
	memDB.sm.Delete(fqn)
}

func (memDB *inMemoryDB) set(fqn string, cake *v1.Cake, op string) {
	memDB.sm.Store(fqn, cake)
}

func (memDB *inMemoryDB) Get(fqn string) *v1.Cake {
	obj, _ := memDB.sm.Load(fqn)
	if obj == nil {
		return nil
	}
	return obj.(*v1.Cake)
}
