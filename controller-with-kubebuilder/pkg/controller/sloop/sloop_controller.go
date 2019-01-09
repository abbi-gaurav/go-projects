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
	"context"
	"fmt"
	"reflect"

	"github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/db"
	shipsv1beta1 "github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/pkg/apis/ships/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")
var database = db.New("memory")

const finalizerName = "sloop.finalizer.ships.com"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Sloop Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSloop{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sloop-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Sloop
	err = c.Watch(&source.Kind{Type: &shipsv1beta1.Sloop{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSloop{}

// ReconcileSloop reconciles a Sloop object
type ReconcileSloop struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Sloop object and makes changes based on the state read
// and what is in the Sloop.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=ships.gaurav.io,resources=sloops,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ships.gaurav.io,resources=sloops/status,verbs=get;update;patch
func (r *ReconcileSloop) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Sloop instance
	instance := &shipsv1beta1.Sloop{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
		}

		if err := r.Update(context.Background(), instance); err != nil {
			return reconcile.Result{Requeue: true}, nil
		}
	} else {
		if containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			if err := r.deleteExternalDependency(instance); err != nil {
				return reconcile.Result{}, err
			}
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
		return reconcile.Result{}, nil
	}

	fqn, err := cache.MetaNamespaceKeyFunc(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	dbObj := database.Get(fqn)
	if dbObj != nil {
		err = r.updateExternal(instance, dbObj, fqn)
	} else {
		err = r.addExternal(instance, fqn)
	}

	return reconcile.Result{}, err
}

func (r *ReconcileSloop) addExternal(instance *shipsv1beta1.Sloop, fqn string) error {
	database.Add(fqn, &instance.Spec)
	instance.Status = shipsv1beta1.SloopStatus{
		Configured: true,
		Update:     false,
	}
	err := r.Status().Update(context.TODO(), instance)
	return err
}

func (r *ReconcileSloop) deleteExternalDependency(instance *shipsv1beta1.Sloop) error {

	fqn, err := cache.MetaNamespaceKeyFunc(instance)
	if err != nil {
		return err
	}

	database.Delete(fqn)
	return nil
}

func (r *ReconcileSloop) updateExternal(instance *shipsv1beta1.Sloop, found *shipsv1beta1.SloopSpec, fqn string) error {
	if !reflect.DeepEqual(&instance.Spec, found) {
		fmt.Println("updating object instance in db", instance.Spec, found)
		newDBObj := instance.Spec
		database.Update(fqn, &newDBObj)
		instance.Status = shipsv1beta1.SloopStatus{
			Configured: true,
			Update:     true,
		}
		err := r.Status().Update(context.TODO(), instance)
		return err
	}

	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
