package prometheus

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
)

// PrometheusOptions contains configuration for Prometheus deployment
type PrometheusOptions struct {
    Name            string
    Namespace       string
    Labels          map[string]string
    // Adding new fields for more configuration flexibility
    ScrapeInterval  string // e.g., "15s"
    RetentionPeriod string // e.g., "15d"
    StorageSize     string // e.g., "10Gi"
}

// Component interface defines methods that each Prometheus component should implement
type Component interface {
    StatefulSet() *appsv1.StatefulSet
    Service() *corev1.Service
}