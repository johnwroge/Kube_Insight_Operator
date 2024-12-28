package prometheus

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Prometheus struct {
	opts PrometheusOptions
}

func New(opts PrometheusOptions) *Prometheus {
	return &Prometheus{
		opts: opts,
	}
}

func (p *Prometheus) StatefulSet() *appsv1.StatefulSet {
	sts := &appsv1.StatefulSet{
		// ... existing ObjectMeta and top-level Spec ...
		Spec: appsv1.StatefulSetSpec{
			// ... existing fields ...
			Template: corev1.PodTemplateSpec{
				// ... existing metadata ...
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:v2.45.0",
							Args: []string{
								"--config.file=/etc/prometheus/prometheus.yml",
								"--storage.tsdb.path=/prometheus",
								fmt.Sprintf("--storage.tsdb.retention.time=%s", p.opts.RetentionPeriod),
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
										Name: p.opts.Name + "-config",
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
								corev1.ResourceStorage: resource.MustParse(p.opts.StorageSize),
							},
						},
					},
				},
			},
		},
	}
	return sts
}

func (p *Prometheus) Service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.opts.Name,
			Namespace: p.opts.Namespace,
			Labels:    p.opts.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "web",
					Port:     9090,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: p.opts.Labels,
		},
	}
}
