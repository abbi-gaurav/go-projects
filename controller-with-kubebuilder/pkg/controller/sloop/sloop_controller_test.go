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
	"fmt"
	shipsv1beta1 "github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

type instanceWithError struct {
	obj *shipsv1beta1.Sloop
	err error
}

type shouldRetry struct {
	flag bool
	err  error
}

var c client.Client

var depKey = types.NamespacedName{Name: "foo", Namespace: "default"}
var expectedRequest = reconcile.Request{NamespacedName: depKey}

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &shipsv1beta1.Sloop{
		ObjectMeta: metav1.ObjectMeta{Name: depKey.Name, Namespace: depKey.Namespace},
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

	fqn, _ := cache.MetaNamespaceKeyFunc(instance)

	create(g, instance, requests, fqn)

	g.Eventually(get(depKey, g).Finalizers).ShouldNot(gomega.BeNil())

	update(depKey, g, "updated", requests, fqn)

	remove(fqn, instance, g, requests)

}

func create(g *gomega.GomegaWithT, instance *shipsv1beta1.Sloop, requests chan reconcile.Request, fqn string) {
	err := c.Create(context.TODO(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	g.Eventually(database.Get(fqn)).Should(gomega.Equal(&instance.Spec))
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
	g.Eventually(func() shouldRetry { return doUpdate(key, g, newRig, requests) }).Should(gomega.Equal(shouldRetry{flag: false, err: nil}))
	g.Eventually(func() string { return database.Get(fqn).Rig }).Should(gomega.Equal(newRig))
}

func doUpdate(key client.ObjectKey, g *gomega.GomegaWithT, newRig string, requests chan reconcile.Request) shouldRetry {
	obj := get(key, g)
	obj.Spec.Rig = newRig
	err := c.Update(context.TODO(), obj)
	g.Eventually(requests).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	if err != nil {
		println(err)
		if errors.IsConflict(err) {
			return shouldRetry{flag: true, err: nil}
		} else {
			return shouldRetry{flag: true, err: err}
		}
	}
	return shouldRetry{flag: false, err: nil}
}

func remove(fqn string, instance *shipsv1beta1.Sloop, g *gomega.GomegaWithT, requests chan reconcile.Request) {
	err := c.Delete(context.TODO(), instance, client.GracePeriodSeconds(0))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	g.Eventually(func() shouldRetry { return verifyRemove(depKey) }).Should(gomega.Equal(shouldRetry{flag: false, err: nil}))

	g.Eventually(func() *shipsv1beta1.SloopSpec { return database.Get(fqn) }).Should(gomega.BeNil())

}

func verifyRemove(key client.ObjectKey) shouldRetry {
	obj := doGet(key)

	if obj.err != nil {
		if errors.IsNotFound(obj.err) {
			return shouldRetry{flag: false, err: nil}
		} else {
			return shouldRetry{flag: true, err: obj.err}
		}
	} else if len(obj.obj.Finalizers) == 0 {
		fmt.Printf("%+v\n", obj.obj)
		return shouldRetry{flag: false, err: nil}
	} else {
		fmt.Printf("%+v\n", obj.obj)
		return shouldRetry{flag: true}
	}
}
