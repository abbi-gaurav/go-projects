/*
My awesome controller
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/clientset/versioned/typed/awesome.controller.io/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAwesomeV1 struct {
	*testing.Fake
}

func (c *FakeAwesomeV1) Cakes(namespace string) v1.CakeInterface {
	return &FakeCakes{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAwesomeV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}