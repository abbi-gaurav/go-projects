/*
My awesome controller
*/
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	awesomecontrolleriov1 "github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	versioned "github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/abbi-gaurav/go-projects/my-awesome-controller/pkg/client/listers/awesome.controller.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CakeInformer provides access to a shared informer and lister for
// Cakes.
type CakeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CakeLister
}

type cakeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCakeInformer constructs a new informer for Cake type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCakeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCakeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCakeInformer constructs a new informer for Cake type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCakeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AwesomeV1().Cakes(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AwesomeV1().Cakes(namespace).Watch(options)
			},
		},
		&awesomecontrolleriov1.Cake{},
		resyncPeriod,
		indexers,
	)
}

func (f *cakeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCakeInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cakeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&awesomecontrolleriov1.Cake{}, f.defaultInformer)
}

func (f *cakeInformer) Lister() v1.CakeLister {
	return v1.NewCakeLister(f.Informer().GetIndexer())
}
