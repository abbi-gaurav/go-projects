/*
My awesome controller
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CakeLister helps list Cakes.
type CakeLister interface {
	// List lists all Cakes in the indexer.
	List(selector labels.Selector) (ret []*v1.Cake, err error)
	// Cakes returns an object that can list and get Cakes.
	Cakes(namespace string) CakeNamespaceLister
	CakeListerExpansion
}

// cakeLister implements the CakeLister interface.
type cakeLister struct {
	indexer cache.Indexer
}

// NewCakeLister returns a new CakeLister.
func NewCakeLister(indexer cache.Indexer) CakeLister {
	return &cakeLister{indexer: indexer}
}

// List lists all Cakes in the indexer.
func (s *cakeLister) List(selector labels.Selector) (ret []*v1.Cake, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Cake))
	})
	return ret, err
}

// Cakes returns an object that can list and get Cakes.
func (s *cakeLister) Cakes(namespace string) CakeNamespaceLister {
	return cakeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CakeNamespaceLister helps list and get Cakes.
type CakeNamespaceLister interface {
	// List lists all Cakes in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Cake, err error)
	// Get retrieves the Cake from the indexer for a given namespace and name.
	Get(name string) (*v1.Cake, error)
	CakeNamespaceListerExpansion
}

// cakeNamespaceLister implements the CakeNamespaceLister
// interface.
type cakeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Cakes in the indexer for a given namespace.
func (s cakeNamespaceLister) List(selector labels.Selector) (ret []*v1.Cake, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Cake))
	})
	return ret, err
}

// Get retrieves the Cake from the indexer for a given namespace and name.
func (s cakeNamespaceLister) Get(name string) (*v1.Cake, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("cake"), name)
	}
	return obj.(*v1.Cake), nil
}
