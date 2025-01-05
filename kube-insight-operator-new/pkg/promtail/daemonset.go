package promtail

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/utils/pointer"
)

// GenerateDaemonSet creates a DaemonSet for Promtail
func (g *ConfigGenerator) GenerateDaemonSet() *appsv1.DaemonSet {
    return &appsv1.DaemonSet{
        ObjectMeta: metav1.ObjectMeta{
            Name:      g.options.Name,
            Namespace: g.options.Namespace,
            Labels:    g.options.Labels,
        },
        Spec: appsv1.DaemonSetSpec{
            Selector: &metav1.LabelSelector{
                MatchLabels: g.options.Labels,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: g.options.Labels,
                },
                Spec: corev1.PodSpec{
                    Tolerations:    g.options.Tolerations,
                    SecurityContext: &corev1.PodSecurityContext{
                        RunAsGroup: pointer.Int64(0),
                        RunAsUser:  pointer.Int64(0),
                    },
                    Containers: []corev1.Container{
                        {
                            Name:  "promtail",
                            Image: "grafana/promtail:2.8.4",
                            Args: []string{
                                "-config.file=/etc/promtail/promtail.yaml",
                                "-client.external-labels=cluster=$(CLUSTER_NAME)",
                            },
                            Env: []corev1.EnvVar{
                                {
                                    Name: "HOSTNAME",
                                    ValueFrom: &corev1.EnvVarSource{
                                        FieldRef: &corev1.ObjectFieldSelector{
                                            FieldPath: "spec.nodeName",
                                        },
                                    },
                                },
                                {
                                    Name: "CLUSTER_NAME",
                                    ValueFrom: &corev1.EnvVarSource{
                                        FieldRef: &corev1.ObjectFieldSelector{
                                            FieldPath: "metadata.namespace",
                                        },
                                    },
                                },
                            },
                            Ports: []corev1.ContainerPort{
                                {
                                    ContainerPort: 9080,
                                    Name:         "http-metrics",
                                },
                            },
                            Resources: *g.options.Resources,
                            SecurityContext: &corev1.SecurityContext{
                                ReadOnlyRootFilesystem: pointer.Bool(true),
                                Capabilities: &corev1.Capabilities{
                                    Drop: []corev1.Capability{"ALL"},
                                },
                            },
                            ReadinessProbe: &corev1.Probe{
                                ProbeHandler: corev1.ProbeHandler{
                                    HTTPGet: &corev1.HTTPGetAction{
                                        Path: "/ready",
                                        Port: intstr.FromInt(9080),
                                    },
                                },
                                InitialDelaySeconds: 10,
                                TimeoutSeconds:      1,
                                PeriodSeconds:      10,
                                SuccessThreshold:    1,
                                FailureThreshold:    5,
                            },
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "config",
                                    MountPath: "/etc/promtail",
                                },
                                {
                                    Name:      "run",
                                    MountPath: "/run/promtail",
                                },
                                {
                                    Name:      "pods",
                                    MountPath: "/var/log/pods",
                                    ReadOnly:  true,
                                },
                                {
                                    Name:      "docker",
                                    MountPath: "/var/lib/docker/containers",
                                    ReadOnly:  true,
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
                                        Name: g.options.Name + "-config",
                                    },
                                },
                            },
                        },
                        {
                            Name: "run",
                            VolumeSource: corev1.VolumeSource{
                                EmptyDir: &corev1.EmptyDirVolumeSource{},
                            },
                        },
                        {
                            Name: "pods",
                            VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                    Path: "/var/log/pods",
                                },
                            },
                        },
                        {
                            Name: "docker",
                            VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                    Path: "/var/lib/docker/containers",
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}