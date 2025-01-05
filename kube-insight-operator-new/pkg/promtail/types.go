package promtail

import (
    corev1 "k8s.io/api/core/v1"
)

type Options struct {
    Name                string
    Namespace           string
    Labels              map[string]string
    LokiURL             string
    Tolerations         []corev1.Toleration
    Resources           *corev1.ResourceRequirements
    ExtraArgs           []string
    ScrapeKubernetesLogs bool
}