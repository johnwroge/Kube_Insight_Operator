/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/johnwroge/kube-insight-operator/api/v1alpha1"
)

// ObservabilityStackReconciler reconciles a ObservabilityStack object
type ObservabilityStackReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile handles the main reconciliation loop for ObservabilityStack
func (r *ObservabilityStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ObservabilityStack instance
	stack := &monitoringv1alpha1.ObservabilityStack{}
	err := r.Get(ctx, req.NamespacedName, stack)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check if Prometheus is enabled and reconcile it
	if stack.Spec.Prometheus.Enabled {
		if err := r.reconcilePrometheus(ctx, stack); err != nil {
			log.Error(err, "Failed to reconcile Prometheus")
			return ctrl.Result{}, err
		}
	}

	// Check if Grafana is enabled and reconcile it
	if stack.Spec.Grafana.Enabled {
		if err := r.reconcileGrafana(ctx, stack); err != nil {
			log.Error(err, "Failed to reconcile Grafana")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// reconcilePrometheus handles the reconciliation of Prometheus resources
func (r *ObservabilityStackReconciler) reconcilePrometheus(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	// Create StatefulSet for Prometheus
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-prometheus", stack.Name),
			Namespace: stack.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "prometheus",
				"app.kubernetes.io/instance":   stack.Name,
				"app.kubernetes.io/managed-by": "kube-insight-operator",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: fmt.Sprintf("%s-prometheus", stack.Name),
			Replicas:    pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     "prometheus",
					"app.kubernetes.io/instance": stack.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":     "prometheus",
						"app.kubernetes.io/instance": stack.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:v2.45.0",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9090,
									Name:          "web",
								},
							},
							// We'll add volume mounts and other configuration later
						},
					},
				},
			},
		},
	}

	// Set controller reference for garbage collection
	if err := ctrl.SetControllerReference(stack, sts, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	// Create or update the StatefulSet
	err := r.Create(ctx, sts)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create Prometheus StatefulSet: %w", err)
	}

	// Create Service for Prometheus
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-prometheus", stack.Name),
			Namespace: stack.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "prometheus",
				"app.kubernetes.io/instance":   stack.Name,
				"app.kubernetes.io/managed-by": "kube-insight-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "web",
					Port:     9090,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name":     "prometheus",
				"app.kubernetes.io/instance": stack.Name,
			},
		},
	}

	// Set controller reference for the service
	if err := ctrl.SetControllerReference(stack, svc, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference for service: %w", err)
	}

	// Create or update the Service
	err = r.Create(ctx, svc)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create Prometheus Service: %w", err)
	}

	return nil
}

// reconcileGrafana handles the reconciliation of Grafana resources
func (r *ObservabilityStackReconciler) reconcileGrafana(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	// TODO: Implement Grafana reconciliation
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObservabilityStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ObservabilityStack{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
