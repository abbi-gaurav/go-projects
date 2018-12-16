package app_test

import (
	"context"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/app"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/clientset/versioned/fake"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions"
	v12 "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	"github.com/satori/go.uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var application *app.Application
var client *fake.Clientset

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	options := opts.DefaultOptions
	informer := newFakeInformer(ctx)
	application = app.New(&options, informer)
	httptest.NewServer(application.ServerMux)
	err := application.Run(ctx.Done())

	println("Set up done...", err)

	retCode := m.Run()
	os.Exit(retCode)
}

func TestCreate(t *testing.T) {
	cake, _ := createNewCake("black-forest", "bavaria", "choclate")
	time.Sleep(10 * time.Millisecond)

	fqn, _ := cache.MetaNamespaceKeyFunc(cake)
	cakeInOven := application.Database.Get(fqn)

	verify(cake, cakeInOven, t)
}

func verify(expected *v1.Cake, actual *v1.Cake, t *testing.T) {
	if expected.Name != actual.Name {
		t.Errorf("name is wrong: want %s, have %s", expected.Name, actual.Name)
	}

	if expected.Spec.Type != actual.Spec.Type {
		t.Errorf("type is wrong: want %s, have %s", expected.Spec.Type, actual.Spec.Type)
	}
}

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
