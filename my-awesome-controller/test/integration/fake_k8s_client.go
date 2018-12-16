package integration

import (
	"context"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/clientset/versioned/fake"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions"
	v12 "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	"github.com/satori/go.uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var client *fake.Clientset

func newFakeInformer(ctx context.Context) v12.CakeInformer {
	client = fake.NewSimpleClientset()
	informerFactory := externalversions.NewSharedInformerFactory(client, 0)
	informer := informerFactory.Awesome().V1().Cakes()

	return informer
}

func createNewCake(name string, namespace string, cakeType string) (*v1.Cake, error) {
	return client.Awesome().Cakes(namespace).Create(getResource(name, namespace, cakeType))
}

func getResource(name string, namespace string, cakeType string) *v1.Cake {
	return &v1.Cake{
		TypeMeta: metav1.TypeMeta{APIVersion: v1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(uuid.NewV4().String()),
		},
		Spec: v1.CakeSpec{
			Type: cakeType,
		},
	}
}
