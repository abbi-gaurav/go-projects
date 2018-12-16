package informer

import (
	clientset "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/clientset/versioned"
	informers "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

func CreateInformer(duration time.Duration) (v1.CakeInformer, clientset.Interface) {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Panicf("error in getting cluster config - %v", err)
	}

	client, err := clientset.NewForConfig(config)

	if err != nil {
		log.Panicf("error in creating client - %v", err)
	}

	factory := informers.NewSharedInformerFactory(client, duration)

	return factory.Awesome().V1().Cakes(), client
}
