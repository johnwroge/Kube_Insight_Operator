package prometheus

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    rbacv1 "k8s.io/api/rbac/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

func NewKubeStateMetricsDeployment(namespace string) *appsv1.Deployment {
    return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "kube-state-metrics",
            Namespace: namespace,
            Labels: map[string]string{
                "app.kubernetes.io/name": "kube-state-metrics",
            },
        },
        Spec: appsv1.DeploymentSpec{
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app.kubernetes.io/name": "kube-state-metrics",
                },
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app.kubernetes.io/name": "kube-state-metrics",
                    },
                },
                Spec: corev1.PodSpec{
                    ServiceAccountName: "kube-state-metrics",
                    Containers: []corev1.Container{
                        {
                            Name:  "kube-state-metrics",
                            Image: "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.10.1",
                            Ports: []corev1.ContainerPort{
                                {
                                    Name:          "metrics",
                                    ContainerPort: 8080,
                                    Protocol:      corev1.ProtocolTCP,
                                },
                                {
                                    Name:          "telemetry",
                                    ContainerPort: 8081,
                                    Protocol:      corev1.ProtocolTCP,
                                },
                            },
                            LivenessProbe: &corev1.Probe{
                                ProbeHandler: corev1.ProbeHandler{
                                    HTTPGet: &corev1.HTTPGetAction{
                                        Path: "/healthz",
                                        Port: intstr.FromInt(8080),
                                    },
                                },
                                InitialDelaySeconds: 5,
                                TimeoutSeconds:      5,
                            },
                            ReadinessProbe: &corev1.Probe{
                                ProbeHandler: corev1.ProbeHandler{
                                    HTTPGet: &corev1.HTTPGetAction{
                                        Path: "/healthz",
                                        Port: intstr.FromInt(8080),
                                    },
                                },
                                InitialDelaySeconds: 5,
                                TimeoutSeconds:      5,
                            },
                        },
                    },
                },
            },
        },
    }
}

func NewKubeStateMetricsService(namespace string) *corev1.Service {
    return &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "kube-state-metrics",
            Namespace: namespace,
            Labels: map[string]string{
                "app.kubernetes.io/name": "kube-state-metrics",
            },
        },
        Spec: corev1.ServiceSpec{
            Ports: []corev1.ServicePort{
                {
                    Name:       "metrics",
                    Port:       8080,
                    TargetPort: intstr.FromString("metrics"),
                    Protocol:   corev1.ProtocolTCP,
                },
                {
                    Name:       "telemetry",
                    Port:       8081,
                    TargetPort: intstr.FromString("telemetry"),
                    Protocol:   corev1.ProtocolTCP,
                },
            },
            Selector: map[string]string{
                "app.kubernetes.io/name": "kube-state-metrics",
            },
        },
    }
}

// RBAC resources for kube-state-metrics
func NewKubeStateMetricsServiceAccount(namespace string) *corev1.ServiceAccount {
    return &corev1.ServiceAccount{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "kube-state-metrics",
            Namespace: namespace,
        },
    }
}

func NewKubeStateMetricsClusterRole() *rbacv1.ClusterRole {
    return &rbacv1.ClusterRole{
        ObjectMeta: metav1.ObjectMeta{
            Name: "kube-state-metrics",
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
                APIGroups: []string{"autoscaling"},
                Resources: []string{
                    "horizontalpodautoscalers",
                },
                Verbs: []string{"list", "watch"},
            },
        },
    }
}

func NewKubeStateMetricsClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
    return &rbacv1.ClusterRoleBinding{
        ObjectMeta: metav1.ObjectMeta{
            Name: "kube-state-metrics",
        },
        RoleRef: rbacv1.RoleRef{
            APIGroup: "rbac.authorization.k8s.io",
            Kind:     "ClusterRole",
            Name:     "kube-state-metrics",
        },
        Subjects: []rbacv1.Subject{
            {
                Kind:      "ServiceAccount",
                Name:      "kube-state-metrics",
                Namespace: namespace,
            },
        },
    }
}