package tempo

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func (g *ConfigGenerator) GenerateStatefulSet() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.options.Name,
			Namespace: g.options.Namespace,
			Labels:    g.options.Labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: g.options.Name,
			Replicas:    pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: g.options.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: g.options.Labels,
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: pointer.Int64(10001),
					},
					Containers: []corev1.Container{
						{
							Name:  "tempo",
							Image: "grafana/tempo:2.3.1",
							Args: []string{
								"-config.file=/etc/tempo/tempo.yaml",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 3200,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "grpc",
									ContainerPort: 9095,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "jaeger-grpc",
									ContainerPort: 14250,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "jaeger-http",
									ContainerPort: 14268,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "otlp-grpc",
									ContainerPort: 4317,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "otlp-http",
									ContainerPort: 4318,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/tempo",
								},
								{
									Name:      "storage",
									MountPath: "/var/tempo",
								},
							},
							Resources: *g.options.Resources,
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.FromInt(3200),
									},
								},
								InitialDelaySeconds: 30,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.FromInt(3200),
									},
								},
								InitialDelaySeconds: 15,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
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
								corev1.ResourceStorage: resource.MustParse(g.options.Storage),
							},
						},
					},
				},
			},
		},
	}
}

func (g *ConfigGenerator) GenerateService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.options.Name,
			Namespace: g.options.Namespace,
			Labels:    g.options.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       3200,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("http"),
				},
				{
					Name:       "grpc",
					Port:       9095,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("grpc"),
				},
				{
					Name:       "jaeger-grpc",
					Port:       14250,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("jaeger-grpc"),
				},
				{
					Name:       "jaeger-http",
					Port:       14268,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("jaeger-http"),
				},
				{
					Name:       "otlp-grpc",
					Port:       4317,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("otlp-grpc"),
				},
				{
					Name:       "otlp-http",
					Port:       4318,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("otlp-http"),
				},
			},
			Selector: g.options.Labels,
		},
	}
}
