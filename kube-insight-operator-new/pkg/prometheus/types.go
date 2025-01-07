package prometheus

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type PrometheusOptions struct {
	Name      string
	Namespace string
	Labels    map[string]string
	ScrapeInterval  string 
	RetentionPeriod string 
	StorageSize     string 
}

type Component interface {
	StatefulSet() *appsv1.StatefulSet
	Service() *corev1.Service
}
