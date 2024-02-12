package controller

import (
	"context"
	"fmt"
	"reflect"

	canaryv1alpha1 "github.com/dishavirk/canary-k8s-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CanaryReconciler reconciles a Canary object
type CanaryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.thefoosthebars.com,resources=canaries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.thefoosthebars.com,resources=canaries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.thefoosthebars.com,resources=canaries/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *CanaryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var canary canaryv1alpha1.Canary
	if err := r.Get(ctx, req.NamespacedName, &canary); err != nil {
		log.Log.Error(err, "unable to fetch Canary")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define a new Deployment object
	canaryDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-canary", canary.Spec.DeploymentName),
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &canary.Spec.Replicas, // Define how many replicas for the canary
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": canary.Spec.DeploymentName, "canary": "true"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": canary.Spec.DeploymentName, "canary": "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: canary.Spec.Image,
						},
					},
				},
			},
		},
	}

	// Set Canary instance as the owner and controller
	if err := controllerutil.SetControllerReference(&canary, canaryDeployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: canaryDeployment.Name, Namespace: canaryDeployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Log.Info("Creating a new Deployment", "Deployment.Namespace", canaryDeployment.Namespace, "Deployment.Name", canaryDeployment.Name,
			"Deployment.NoOfReplicas", canaryDeployment.Spec.Replicas)
		err = r.Create(ctx, canaryDeployment)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update the Canary status with the pod names
	// List the pods for this canary's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(canaryDeployment.Namespace),
		client.MatchingLabels(labelsForCanary(canary.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Log.Error(err, "Failed to list pods", "Canary.Namespace", canary.Namespace, "Canary.Name", canary.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	log.Log.Info("Canary deployment pods", "Pod Names", podNames)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, canary.Status.Nodes) {
		canary.Status.Nodes = podNames
		err := r.Status().Update(ctx, &canary)
		if err != nil {
			log.Log.Error(err, "Failed to update Canary status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CanaryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&canaryv1alpha1.Canary{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

// labelsForCanary returns the labels for selecting the resources
// belonging to the given canary CR name.
func labelsForCanary(name string) map[string]string {
	return map[string]string{"type": "canary", "cr_name": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
