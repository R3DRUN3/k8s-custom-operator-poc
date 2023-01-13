package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	tutorialv1 "test.domain/poc/api/v1"
)

// MyCustomResourceReconciler reconciles a MyCustomResource object
type MyCustomResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC permissions to monitor MyCustomResource custom resources
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=MyCustomResources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=MyCustomResources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=MyCustomResources/finalizers,verbs=update

// RBAC permissions to monitor pods
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MyCustomResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling MyCustomResource custom resource")

	// Get the MyCustomResource resource that triggered the reconciliation request
	var MyCustomResource tutorialv1.MyCustomResource
	if err := r.Get(ctx, req.NamespacedName, &MyCustomResource); err != nil {
		log.Error(err, "unable to fetch MyCustomResource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get pods with the same name as MyCustomResource's friend
	var podList corev1.PodList
	var friendFound bool
	if err := r.List(ctx, &podList); err != nil {
		log.Error(err, "unable to list pods")
	} else {
		for _, item := range podList.Items {
			if item.GetName() == MyCustomResource.Spec.Name {
				log.Info("pod linked to a MyCustomResource custom resource found", "name", item.GetName())
				friendFound = true
			}
		}
	}

	// Update MyCustomResource' Healthy status
	MyCustomResource.Status.Healthy = friendFound
	if err := r.Status().Update(ctx, &MyCustomResource); err != nil {
		log.Error(err, "unable to update MyCustomResource's Healthy status", "status", friendFound)
		return ctrl.Result{}, err
	}
	log.Info("MyCustomResource's Healthy status updated", "status", friendFound)

	log.Info("MyCustomResource custom resource reconciled")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyCustomResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tutorialv1.MyCustomResource{}).
		Watches(
			&source.Kind{Type: &corev1.Pod{}},
			handler.EnqueueRequestsFromMapFunc(r.mapPodsReqToMyCustomResourceReq),
		).
		Complete(r)
}

func (r *MyCustomResourceReconciler) mapPodsReqToMyCustomResourceReq(obj client.Object) []reconcile.Request {
	ctx := context.Background()
	log := log.FromContext(ctx)

	// List all the MyCustomResource custom resource
	req := []reconcile.Request{}
	var list tutorialv1.MyCustomResourceList
	if err := r.Client.List(context.TODO(), &list); err != nil {
		log.Error(err, "unable to list MyCustomResource custom resources")
	} else {
		// Only keep MyCustomResource custom resources related to the Pod that triggered the reconciliation request
		for _, item := range list.Items {
			if item.Spec.Name == obj.GetName() {
				req = append(req, reconcile.Request{
					NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace},
				})
				log.Info("pod linked to a MyCustomResource custom resource issued an event", "name", obj.GetName())
			}
		}
	}
	return req
}
