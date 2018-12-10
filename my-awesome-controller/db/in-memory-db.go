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

func (memDB *inMemoryDB) Add(cake *v1.Cake) {
	memDB.set(cake, "add")
}

func (memDB *inMemoryDB) Update(cake *v1.Cake) {
	memDB.set(cake, "update")
}

func (memDB *inMemoryDB) Delete(cake *v1.Cake) {
	memDB.Lock()
	delete(memDB.m, fqName(cake))
	memDB.log("delete")
	memDB.Unlock()
}

func (memDB *inMemoryDB) set(cake *v1.Cake, op string) {
	memDB.Lock()
	memDB.m[cake.Name] = cake
	memDB.log(op)
	memDB.Unlock()
}

func (memDB *inMemoryDB) log(op string) {
	fmt.Printf("Map post operation %s is %v", op, memDB.m)
}
