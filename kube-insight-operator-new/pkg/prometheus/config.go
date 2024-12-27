package prometheus

import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigGenerator handles Prometheus configuration
type ConfigGenerator struct {
    Name      string
    Namespace string
    Labels    map[string]string
}

func NewConfigGenerator(name, namespace string, labels map[string]string) *ConfigGenerator {
    return &ConfigGenerator{
        Name:      name,
        Namespace: namespace,
        Labels:    labels,
    }
}

func (c *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      c.Name + "-config",
            Namespace: c.Namespace,
            Labels:    c.Labels,
        },
        Data: map[string]string{
            "prometheus.yml": `
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']
  
  - job_name: 'kube-state-metrics'
    static_configs:
      - targets: ['kube-state-metrics:8080']`,
        },
    }
}