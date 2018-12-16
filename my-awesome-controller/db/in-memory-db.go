package db

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"sync"
)

type inMemoryDB struct {
	sync.RWMutex
	m map[string]*v1.Cake
}

func newInMemory() DB {
	return &inMemoryDB{
		m: make(map[string]*v1.Cake),
	}
}

func (memDB *inMemoryDB) Add(fqn string, cake *v1.Cake) {
	memDB.set(fqn, cake, "add")
}

func (memDB *inMemoryDB) Update(fqn string, cake *v1.Cake) {
	memDB.set(fqn, cake, "update")
}

func (memDB *inMemoryDB) Delete(fqn string) {
	memDB.Lock()
	delete(memDB.m, fqn)
	memDB.log("delete")
	memDB.Unlock()
}

func (memDB *inMemoryDB) set(fqn string, cake *v1.Cake, op string) {
	memDB.Lock()
	memDB.m[fqn] = cake
	memDB.log(op)
	memDB.Unlock()
}

func (memDB *inMemoryDB) Get(fqn string) *v1.Cake {
	memDB.RLock()
	obj := memDB.m[fqn]
	return obj
}

func (memDB *inMemoryDB) log(op string) {
	fmt.Printf("Map post operation %s is %v", op, memDB.m)
}
