package label

import (
	"context"

	labv1 "github.com/lab/labels-operator/pkg/apis/lab/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_label")

// Add creates a new Label Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLabel{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("label-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Label
	err = c.Watch(&source.Kind{Type: &labv1.Label{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to the pods from all namespaces
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLabel implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLabel{}

// ReconcileLabel reconciles a Label object
type ReconcileLabel struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Label object and makes changes based on the state read
// and what is in the Label.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLabel) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Labels")

	instance := &labv1.Label{}

	r.client.Get(context.TODO(), client.ObjectKey{Namespace: "labels-operator", Name: "labels-operator"}, instance)

	// go through all projects defined on the CRD
	for _, namespace := range instance.Spec.Projects {
		podList := &corev1.PodList{}
		err := r.client.List(context.TODO(), podList, client.InNamespace(namespace))
		if err != nil {
			return reconcile.Result{}, err
		}

		// add the label to all pods
		for _, pod := range podList.Items {
			found := &corev1.Pod{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.ObjectMeta.Name, Namespace: pod.ObjectMeta.Namespace}, found)

			if err != nil {
				return reconcile.Result{}, err
			}

			// declare the pod.ObjectMeta.Labels for pods that do not had labels previously defined
			if len(pod.ObjectMeta.Labels) == 0 {
				pod.ObjectMeta.Labels = make(map[string]string)
			}

			pod.ObjectMeta.Labels["created-by"] = "labels-operator"
			pod.ObjectMeta.Labels["pod-namespace"] = pod.ObjectMeta.Namespace

			// add the label only to pods in a Running state
			if pod.Status.Phase == "Running" {
				err = r.client.Update(context.TODO(), &pod)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}
	return reconcile.Result{}, nil
}
