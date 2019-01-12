/*
Copyright 2018 Gaurav Abbi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sloop

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
	"testing"
	"time"

	shipsv1beta1 "github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type instanceWithError struct {
	obj *shipsv1beta1.Sloop
	err error
}

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
var depKey = types.NamespacedName{Name: "foo", Namespace: "default"}

const timeout = time.Second * 10

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &shipsv1beta1.Sloop{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
		Spec:       shipsv1beta1.SloopSpec{Rig: "test-rig"},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	err = c.Create(context.TODO(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	fqn, _ := cache.MetaNamespaceKeyFunc(instance)
	g.Expect(database.Get(fqn)).To(gomega.Equal(&instance.Spec))

	g.Eventually(get(depKey, g).Finalizers, timeout).ShouldNot(gomega.BeNil())

	update(depKey, g, "updated", requests, fqn)

	remove(fqn, instance, g, requests)

}

func get(key client.ObjectKey, g *gomega.GomegaWithT) *shipsv1beta1.Sloop {
	obj := doGet(key)
	g.Expect(obj.err).NotTo(gomega.HaveOccurred())
	return obj.obj
}

func doGet(key client.ObjectKey) instanceWithError {
	obj := &shipsv1beta1.Sloop{}
	err := c.Get(context.TODO(), key, obj)
	return instanceWithError{
		obj: obj,
		err: err,
	}
}

func update(key client.ObjectKey, g *gomega.GomegaWithT, newRig string, requests chan reconcile.Request, fqn string) {
	time.Sleep(10 * time.Second)
	err := doUpdate(key, g, newRig, requests)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	time.Sleep(500 * time.Millisecond)
	g.Eventually(database.Get(fqn).Rig, timeout).Should(gomega.Equal(newRig))
}

func doUpdate(key client.ObjectKey, g *gomega.GomegaWithT, newRig string, requests chan reconcile.Request) error {
	obj := get(key, g)
	obj.Spec.Rig = newRig
	err := c.Update(context.TODO(), obj)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	return err
}

func remove(fqn string, instance *shipsv1beta1.Sloop, g *gomega.GomegaWithT, requests chan reconcile.Request) {
	err := c.Delete(context.TODO(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	time.Sleep(500 * time.Millisecond)
	g.Eventually(verifyDelete(depKey), timeout).Should(gomega.Equal(true))

	g.Eventually(database.Get(fqn), timeout).Should(gomega.BeNil())
}

func verifyDelete(key client.ObjectKey) bool {
	obj := doGet(key)
	if obj.err != nil && errors.IsNotFound(obj.err) {
		return true
	}

	if len(obj.obj.Finalizers) == 0 {
		return true
	}

	return false
}
