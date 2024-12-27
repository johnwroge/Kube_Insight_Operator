package prometheus

import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "fmt"
)

// ConfigGenerator handles Prometheus configuration
// type ConfigGenerator struct {
//     Name      string
//     Namespace string
//     Labels    map[string]string
// }

type ConfigGenerator struct {
  Name            string
  Namespace       string
  Labels          map[string]string
  ScrapeInterval  string // Allow customization
  RetentionTime   string // For data retention
}

func NewConfigGenerator(name, namespace string, labels map[string]string) *ConfigGenerator {
  if labels == nil {
      labels = make(map[string]string)
  }
  return &ConfigGenerator{
      Name:      name,
      Namespace: namespace,
      Labels:    labels,
  }
}

// func (c *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
//     return &corev1.ConfigMap{
//         ObjectMeta: metav1.ObjectMeta{
//             Name:      c.Name + "-config",
//             Namespace: c.Namespace,
//             Labels:    c.Labels,
//         },
//         Data: map[string]string{
//             "prometheus.yml": `
// global:
//   scrape_interval: 15s
//   evaluation_interval: 15s

// scrape_configs:
//   - job_name: 'prometheus'
//     static_configs:
//       - targets: ['localhost:9090']
  
//   - job_name: 'node-exporter'
//     static_configs:
//       - targets: ['node-exporter:9100']
  
//   - job_name: 'kube-state-metrics'
//     static_configs:
//       - targets: ['kube-state-metrics:8080']`,
//         },
//     }
// }



func (c *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
  scrapeInterval := c.ScrapeInterval
  if scrapeInterval == "" {
      scrapeInterval = "15s"  // Default
  }

  config := fmt.Sprintf(`
global:
scrape_interval: %s
evaluation_interval: %s

scrape_configs:
- job_name: 'prometheus'
  static_configs:
    - targets: ['localhost:9090']

- job_name: 'node-exporter'
  static_configs:
    - targets: ['node-exporter:9100']

- job_name: 'kube-state-metrics'
  static_configs:
    - targets: ['kube-state-metrics:8080']`, scrapeInterval, scrapeInterval)

  return &corev1.ConfigMap{
      ObjectMeta: metav1.ObjectMeta{
          Name:      c.Name + "-config",
          Namespace: c.Namespace,
          Labels:    c.Labels,
      },
      Data: map[string]string{
          "prometheus.yml": config,
      },
  }
}