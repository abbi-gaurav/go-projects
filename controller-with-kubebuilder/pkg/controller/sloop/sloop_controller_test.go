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
	shipsv1beta1 "github.com/abbi-gaurav/go-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
	"testing"
)

type instanceWithError struct {
	obj *shipsv1beta1.Sloop
	err error
}

type testDriver struct {
	manager        manager.Manager
	requestChannel chan reconcile.Request
	stopMgr        chan struct{}
	mgrStopped     *sync.WaitGroup
	g              *gomega.GomegaWithT
	depKey         types.NamespacedName
}

type shouldRetry struct {
	flag bool
	err  error
}

var c client.Client

func TestCreateDelete(t *testing.T) {
	td := testSetUp(t, "test-create-delete")
	g := td.g

	instance := td.getInstanceObj("test-rig")

	defer func() {
		td.close()
	}()

	fqn, _ := cache.MetaNamespaceKeyFunc(instance)

	td.create(instance, fqn)
	g.Eventually(get(td.depKey, g).Finalizers).ShouldNot(gomega.BeNil())
	td.remove(fqn, instance)

}

func TestCreateUpdate(t *testing.T) {
	td := testSetUp(t, "test-create-update")

	defer func() {
		td.close()
	}()

	instance := td.getInstanceObj("test-update")
	fqn, _ := cache.MetaNamespaceKeyFunc(instance)

	td.create(instance, fqn)
	td.update("updated", fqn)
}

func (td *testDriver) close() {
	close(td.stopMgr)
	td.mgrStopped.Wait()
}

func testSetUp(t *testing.T, name string) *testDriver {
	g := gomega.NewGomegaWithT(t)
	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	return &testDriver{
		manager:        mgr,
		requestChannel: requests,
		stopMgr:        stopMgr,
		mgrStopped:     mgrStopped,
		g:              g,
		depKey:         types.NamespacedName{Name: name, Namespace: "default"},
	}

}

func (td *testDriver) getInstanceObj(rig string) *shipsv1beta1.Sloop {
	return &shipsv1beta1.Sloop{
		ObjectMeta: metav1.ObjectMeta{Name: td.depKey.Name, Namespace: td.depKey.Namespace},
		Spec:       shipsv1beta1.SloopSpec{Rig: rig},
	}
}

func (td *testDriver) create(instance *shipsv1beta1.Sloop, fqn string) {
	g := td.g
	requests := td.requestChannel
	err := c.Create(context.Background(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests).Should(gomega.Receive(gomega.Equal(reconcile.Request{NamespacedName: td.depKey})))
	g.Eventually(database.Get(fqn)).Should(gomega.Equal(&instance.Spec))
}

func get(key client.ObjectKey, g *gomega.GomegaWithT) *shipsv1beta1.Sloop {
	obj := doGet(key)
	g.Expect(obj.err).NotTo(gomega.HaveOccurred())
	return obj.obj
}

func doGet(key client.ObjectKey) instanceWithError {
	obj := &shipsv1beta1.Sloop{}
	err := c.Get(context.Background(), key, obj)
	return instanceWithError{
		obj: obj,
		err: err,
	}
}

func (td *testDriver) update(newRig string, fqn string) {
	g := td.g
	g.Eventually(func() shouldRetry { return td.doUpdate(newRig) }).Should(gomega.Equal(shouldRetry{flag: false, err: nil}))
	g.Eventually(func() string { return database.Get(fqn).Rig }).Should(gomega.Equal(newRig))
}

func (td *testDriver) doUpdate(newRig string) shouldRetry {
	g := td.g
	key := td.depKey
	requests := td.requestChannel
	expectedRequest := reconcile.Request{NamespacedName: key}
	obj := get(key, g)
	obj.Spec.Rig = newRig
	err := c.Update(context.Background(), obj)
	if err != nil {
		println(err)
		if errors.IsConflict(err) {
			return shouldRetry{flag: true, err: nil}
		} else {
			return shouldRetry{flag: true, err: err}
		}
	}
	g.Eventually(requests).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	return shouldRetry{flag: false, err: nil}
}

func (td *testDriver) remove(fqn string, instance *shipsv1beta1.Sloop) {
	g := td.g
	requests := td.requestChannel
	depKey := td.depKey
	expectedRequest := reconcile.Request{NamespacedName: depKey}
	err := c.Delete(context.Background(), instance, client.GracePeriodSeconds(0))
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
