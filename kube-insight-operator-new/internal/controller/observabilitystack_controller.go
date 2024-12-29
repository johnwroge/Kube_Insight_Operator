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

//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.example.com,resources=observabilitystacks/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

package controller

import (
	"context"
	"fmt"

	monitoringv1alpha1 "github.com/johnwroge/kube-insight-operator/kube-insight-operator-new/api/v1alpha1"
	"github.com/johnwroge/kube-insight-operator/kube-insight-operator-new/pkg/grafana"
	"github.com/johnwroge/kube-insight-operator/kube-insight-operator-new/pkg/prometheus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ObservabilityStackReconciler reconciles a ObservabilityStack object
type ObservabilityStackReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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

func (r *ObservabilityStackReconciler) reconcilePrometheus(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {

	if err := r.reconcilePrometheusRBAC(ctx, stack); err != nil {
		return fmt.Errorf("failed to reconcile Prometheus RBAC: %w", err)
	}

	if err := r.reconcileKubeStateMetrics(ctx, stack); err != nil {
		return fmt.Errorf("failed to reconcile kube-state-metrics: %w", err)
	}

	// Define common labels
	labels := map[string]string{
		"app.kubernetes.io/name":       "prometheus",
		"app.kubernetes.io/instance":   stack.Name,
		"app.kubernetes.io/managed-by": "kube-insight-operator",
	}

	// Create ConfigMap first
	configGen := prometheus.NewConfigGenerator(
		fmt.Sprintf("%s-prometheus", stack.Name),
		stack.Namespace,
		labels,
	)

	configMap := configGen.GenerateConfigMap()
	if err := ctrl.SetControllerReference(stack, configMap, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on configmap: %w", err)
	}

	if err := r.createOrUpdate(ctx, configMap); err != nil {
		return fmt.Errorf("failed to reconcile Prometheus ConfigMap: %w", err)
	}

	// Create StatefulSet for Prometheus
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-prometheus", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
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
					ServiceAccountName:           fmt.Sprintf("%s-prometheus", stack.Name),
					AutomountServiceAccountToken: pointer.Bool(true),
					Containers: []corev1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:v2.45.0",
							Args: []string{
								"--config.file=/etc/prometheus/prometheus.yml",
								"--storage.tsdb.path=/prometheus",
								"--storage.tsdb.retention.time=15d",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9090,
									Name:          "web",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/prometheus",
								},
								{
									Name:      "storage",
									MountPath: "/prometheus",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-prometheus-config", stack.Name),
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "storage",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(stack.Spec.Prometheus.Storage),
							},
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

	// Create or update the StatefulSet using our helper
	if err := r.createOrUpdate(ctx, sts); err != nil {
		return fmt.Errorf("failed to reconcile Prometheus StatefulSet: %w", err)
	}

	// Create Service for Prometheus
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-prometheus", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
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

	// Create or update the Service using our helper
	if err := r.createOrUpdate(ctx, svc); err != nil {
		return fmt.Errorf("failed to reconcile Prometheus Service: %w", err)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObservabilityStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ObservabilityStack{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// Helper function for creating or updating resources
func (r *ObservabilityStackReconciler) createOrUpdate(ctx context.Context, obj client.Object) error {
	err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj.DeepCopyObject().(client.Object))
	if err != nil {
		if errors.IsNotFound(err) {
			if err = r.Create(ctx, obj); err != nil {
				return fmt.Errorf("failed to create resource: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get resource: %w", err)
	}

	if err = r.Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}
	return nil
}

func (r *ObservabilityStackReconciler) reconcilePrometheusRBAC(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	// Create ServiceAccount (your existing code)
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-prometheus", stack.Name),
			Namespace: stack.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":     "prometheus",
				"app.kubernetes.io/instance": stack.Name,
			},
		},
	}

	if err := ctrl.SetControllerReference(stack, sa, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on serviceaccount: %w", err)
	}

	if err := r.createOrUpdate(ctx, sa); err != nil {
		return fmt.Errorf("failed to reconcile ServiceAccount: %w", err)
	}

	// Create ClusterRole
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-prometheus", stack.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "prometheus",
				"app.kubernetes.io/instance": stack.Name,
			},
		},

		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"nodes",
					"nodes/proxy",
					"nodes/metrics",
					"services",
					"endpoints",
					"pods",
				},
				Verbs: []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				NonResourceURLs: []string{
					"/metrics",
					"/api",
					"/api/*",
				},
				Verbs: []string{"get"},
			},
		},
	}

	if err := r.createOrUpdate(ctx, cr); err != nil {
		return fmt.Errorf("failed to reconcile ClusterRole: %w", err)
	}

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-prometheus", stack.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "prometheus",
				"app.kubernetes.io/instance": stack.Name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     fmt.Sprintf("%s-prometheus", stack.Name),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      fmt.Sprintf("%s-prometheus", stack.Name),
				Namespace: stack.Namespace,
			},
		},
	}

	if err := r.createOrUpdate(ctx, crb); err != nil {
		return fmt.Errorf("failed to reconcile ClusterRoleBinding: %w", err)
	}

	return nil
}

func (r *ObservabilityStackReconciler) reconcileKubeStateMetricsRBAC(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	// Skip if kube-state-metrics is not enabled
	if !stack.Spec.Prometheus.KubeStateMetrics.Enabled {
		return nil
	}

	// Common labels
	labels := map[string]string{
		"app.kubernetes.io/name":     "kube-state-metrics",
		"app.kubernetes.io/instance": stack.Name,
	}

	// Create ServiceAccount
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-kube-state-metrics", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
		},
	}

	if err := ctrl.SetControllerReference(stack, sa, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on serviceaccount: %w", err)
	}

	if err := r.createOrUpdate(ctx, sa); err != nil {
		return fmt.Errorf("failed to reconcile ServiceAccount: %w", err)
	}

	// Create ClusterRole
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-kube-state-metrics", stack.Name),
			Labels: labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"configmaps",
					"secrets",
					"nodes",
					"pods",
					"services",
					"resourcequotas",
					"replicationcontrollers",
					"limitranges",
					"persistentvolumeclaims",
					"persistentvolumes",
					"namespaces",
					"endpoints",
				},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{
					"statefulsets",
					"daemonsets",
					"deployments",
					"replicasets",
				},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{
					"cronjobs",
					"jobs",
				},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{
					"storageclasses",
				},
				Verbs: []string{"list", "watch"},
			},
		},
	}

	if err := r.createOrUpdate(ctx, cr); err != nil {
		return fmt.Errorf("failed to reconcile ClusterRole: %w", err)
	}

	// Create ClusterRoleBinding
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-kube-state-metrics", stack.Name),
			Labels: labels,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     cr.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		},
	}

	if err := r.createOrUpdate(ctx, crb); err != nil {
		return fmt.Errorf("failed to reconcile ClusterRoleBinding: %w", err)
	}

	return nil
}

func (r *ObservabilityStackReconciler) reconcileKubeStateMetrics(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	if !stack.Spec.Prometheus.KubeStateMetrics.Enabled {
		return nil
	}

	if err := r.reconcileKubeStateMetricsRBAC(ctx, stack); err != nil {
		return fmt.Errorf("failed to reconcile kube-state-metrics RBAC: %w", err)
	}

	labels := map[string]string{
		"app.kubernetes.io/name":     "kube-state-metrics",
		"app.kubernetes.io/instance": stack.Name,
	}

	// Create deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-kube-state-metrics", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fmt.Sprintf("%s-kube-state-metrics", stack.Name),
					Containers: []corev1.Container{
						{
							Name:  "kube-state-metrics",
							Image: "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.10.0",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http-metrics",
									ContainerPort: 8080,
								},
								{
									Name:          "telemetry",
									ContainerPort: 8081,
								},
							},
						},
					},
				},
			},
		},
	}

	// Set controller reference
	if err := ctrl.SetControllerReference(stack, deployment, r.Scheme); err != nil {
		return err
	}

	// Create or update deployment
	if err := r.createOrUpdate(ctx, deployment); err != nil {
		return err
	}

	// Create service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-kube-state-metrics", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http-metrics",
					Port:       8080,
					TargetPort: intstr.FromString("http-metrics"),
				},
				{
					Name:       "telemetry",
					Port:       8081,
					TargetPort: intstr.FromString("telemetry"),
				},
			},
			Selector: labels,
		},
	}

	// Set controller reference
	if err := ctrl.SetControllerReference(stack, service, r.Scheme); err != nil {
		return err
	}

	// Create or update service
	if err := r.createOrUpdate(ctx, service); err != nil {
		return err
	}

	return nil
}

func (r *ObservabilityStackReconciler) reconcileGrafana(ctx context.Context, stack *monitoringv1alpha1.ObservabilityStack) error {
	if !stack.Spec.Grafana.Enabled {
		return nil
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "grafana",
		"app.kubernetes.io/instance":   stack.Name,
		"app.kubernetes.io/managed-by": "kube-insight-operator",
	}

	// Create Grafana instance
	grafanaOpts := grafana.Options{
		Name:          fmt.Sprintf("%s-grafana", stack.Name),
		Namespace:     stack.Namespace,
		Labels:        labels,
		AdminPassword: stack.Spec.Grafana.AdminPassword,
		Storage:       stack.Spec.Grafana.Storage,
		PrometheusURL: fmt.Sprintf("http://%s-prometheus:9090", stack.Name),
	}

	g := grafana.New(grafanaOpts)

	// Create ConfigMap
	configMap := g.GenerateConfigMap()
	if err := ctrl.SetControllerReference(stack, configMap, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on configmap: %w", err)
	}

	if err := r.createOrUpdate(ctx, configMap); err != nil {
		return fmt.Errorf("failed to reconcile Grafana ConfigMap: %w", err)
	}


	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-grafana", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "init-chown-data",
							Image: "busybox:1.35",
							Command: []string{
								"sh",
								"-c",
								"mkdir -p /etc/grafana/provisioning/datasources /etc/grafana/provisioning/dashboards && chown -R 472:472 /etc/grafana /var/lib/grafana",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/grafana",
								},
								{
									Name:      "storage",
									MountPath: "/var/lib/grafana",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "grafana",
							Image: "grafana/grafana:9.5.3",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 3000,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/grafana/grafana.ini",
									SubPath:   "grafana.ini",
								},
								{
									Name:      "config",
									MountPath: "/etc/grafana/provisioning/datasources/datasources.yaml",
									SubPath:   "datasources.yaml",
								},
								{
									Name:      "storage",
									MountPath: "/var/lib/grafana",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
								},
							},
						},
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: fmt.Sprintf("%s-grafana", stack.Name),
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(stack, deployment, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on deployment: %w", err)
	}

	if err := r.createOrUpdate(ctx, deployment); err != nil {
		return fmt.Errorf("failed to reconcile Grafana Deployment: %w", err)
	}

	// Create Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-grafana", stack.Name),
			Namespace: stack.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     3000,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: labels,
		},
	}

	if err := ctrl.SetControllerReference(stack, svc, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on service: %w", err)
	}

	if err := r.createOrUpdate(ctx, svc); err != nil {
		return fmt.Errorf("failed to reconcile Grafana Service: %w", err)
	}

	// Create PVC for Grafana storage
	if stack.Spec.Grafana.Storage != "" {
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-grafana", stack.Name),
				Namespace: stack.Namespace,
				Labels:    labels,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: pointer.String("standard"),
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(stack.Spec.Grafana.Storage),
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(stack, pvc, r.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference on pvc: %w", err)
		}

		if err := r.createOrUpdate(ctx, pvc); err != nil {
			return fmt.Errorf("failed to reconcile Grafana PVC: %w", err)
		}
	}

	return nil
}
