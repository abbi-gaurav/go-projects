/*

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

package genericdaemon

import (
	"context"
	"reflect"

	mygroupv1beta1 "github.com/abbi-gaurav/go-projects/operator-with-kubebuilder/pkg/apis/mygroup/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new GenericDaemon Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGenericDaemon{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("genericdaemon-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to GenericDaemon
	err = c.Watch(&source.Kind{Type: &mygroupv1beta1.GenericDaemon{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by GenericDaemon - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mygroupv1beta1.GenericDaemon{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileGenericDaemon{}

// ReconcileGenericDaemon reconciles a GenericDaemon object
type ReconcileGenericDaemon struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a GenericDaemon object and makes changes based on the state read
// and what is in the GenericDaemon.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mygroup.gabbi.io,resources=genericdaemons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mygroup.gabbi.io,resources=genericdaemons/status,verbs=get;update;patch
func (r *ReconcileGenericDaemon) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the GenericDaemon instance
	instance := &mygroupv1beta1.GenericDaemon{}
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

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Daemonset object
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-daemonset",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"daemonset": instance.Name + "-daemonset"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"daemonset": instance.Name + "-daemonset"}},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{"daemon": instance.Spec.Label},
					Containers: []corev1.Container{
						{
							Name:  "genericdaemon",
							Image: instance.Spec.Image,
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, daemonset, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the DaemonSet already exists
	found := &appsv1.DaemonSet{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: daemonset.Name, Namespace: daemonset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating DaemonSet", "namespace", daemonset.Namespace, "name", daemonset.Name)
		err = r.Create(context.TODO(), daemonset)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//Get the number of ready daemonsets and set the count status
	if found.Status.NumberReady != instance.Status.Count {
		log.Info("Updating status %s/%s\n", instance.Namespace, instance.Name)
		instance.Status.Count = found.Status.NumberReady
		err = r.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(daemonset.Spec, found.Spec) {
		found.Spec = daemonset.Spec
		log.Info("Updating DaemonSet", "namespace", daemonset.Namespace, "name", daemonset.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
