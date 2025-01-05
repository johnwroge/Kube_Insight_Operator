package tempo

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (g *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
	// Convert retention days to hours
	retentionPeriod := fmt.Sprintf("%dh", g.options.RetentionDays*24)

	tempoConfig := `
server:
  http_listen_port: 3200

distributor:
  receivers:
    jaeger:
      protocols:
        grpc:
          endpoint: "0.0.0.0:14250"
        thrift_http:
          endpoint: "0.0.0.0:14268"
    otlp:
      protocols:
        grpc:
          endpoint: "0.0.0.0:4317"
        http:
          endpoint: "0.0.0.0:4318"

ingester:
  max_block_duration: 5m
  trace_idle_period: 10s

compactor:
  compaction:
    block_retention: ` + retentionPeriod + `

storage:
  trace:
    backend: local
    local:
      path: /var/tempo/traces
    wal:
      path: /var/tempo/wal

overrides:
  defaults:
    metrics_generator:
      processors: [service-graphs, span-metrics]`

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.options.Name + "-config",
			Namespace: g.options.Namespace,
			Labels:    g.options.Labels,
		},
		Data: map[string]string{
			"tempo.yaml": tempoConfig,
		},
	}
}
