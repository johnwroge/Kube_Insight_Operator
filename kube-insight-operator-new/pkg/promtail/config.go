package promtail

import (
    "fmt"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
    Name        string
    Namespace   string
    Labels      map[string]string
    LokiURL     string
}

type ConfigGenerator struct {
    options Options
}

func NewConfigGenerator(opts Options) *ConfigGenerator {
    return &ConfigGenerator{
        options: opts,
    }
}

func (g *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
    promtailConfig := fmt.Sprintf(`
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /run/promtail/positions.yaml

clients:
  - url: %s/loki/api/v1/push

scrape_configs:
  - job_name: kubernetes-pods
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels:
          - __meta_kubernetes_pod_controller_name
        regex: ([0-9a-z-.]+?)(-[0-9a-f]{8,10})?
        action: replace
        target_label: __tmp_controller_name
      - source_labels:
          - __meta_kubernetes_pod_label_app_kubernetes_io_name
          - __meta_kubernetes_pod_label_app
          - __tmp_controller_name
          - __meta_kubernetes_pod_name
        regex: ^;*([^;]+)(;.*)?$
        action: replace
        target_label: app
      - source_labels:
          - __meta_kubernetes_pod_node_name
        action: replace
        target_label: node_name
      - source_labels:
          - __meta_kubernetes_namespace
        action: replace
        target_label: namespace
      - source_labels:
          - __meta_kubernetes_pod_name
        action: replace
        target_label: pod
      - source_labels:
          - __meta_kubernetes_pod_container_name
        action: replace
        target_label: container
      - action: replace
        source_labels:
          - __meta_kubernetes_namespace
        target_label: exported_namespace
      - action: replace
        source_labels:
          - __meta_kubernetes_pod_name
        target_label: exported_pod
      - action: replace
        source_labels:
          - __meta_kubernetes_pod_container_name
        target_label: exported_container
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)
      - replacement: /var/log/pods/*$1/*.log
        separator: /
        source_labels:
          - __meta_kubernetes_pod_uid
          - __meta_kubernetes_pod_container_name
        target_label: __path__`, g.options.LokiURL)

    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      g.options.Name + "-config",
            Namespace: g.options.Namespace,
            Labels:    g.options.Labels,
        },
        Data: map[string]string{
            "promtail.yaml": promtailConfig,
        },
    }
}