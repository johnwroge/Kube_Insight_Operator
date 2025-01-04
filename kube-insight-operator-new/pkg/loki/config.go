package loki

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name          string
	Namespace     string
	Labels        map[string]string
	Storage       string
	RetentionDays int32
}

type ConfigGenerator struct {
	options Options
}

func NewConfigGenerator(opts Options) *ConfigGenerator {
	return &ConfigGenerator{
		options: opts,
	}
}

// func (g *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
// 	lokiConfig := `
// auth_enabled: false

// server:
//   http_listen_port: 3100
//   grpc_listen_port: 9096

// common:
//   path_prefix: /loki
//   storage:
//     filesystem:
//       chunks_directory: /loki/chunks
//       rules_directory: /loki/rules
//   replication_factor: 1
//   ring:
//     kvstore:
//       store: inmemory

// schema_config:
//   configs:
//     - from: 2020-10-24
//       store: boltdb-shipper
//       object_store: filesystem
//       schema: v11
//       index:
//         prefix: index_
//         period: 24h

// storage_config:
//   boltdb_shipper:
//     active_index_directory: /loki/index
//     cache_location: /loki/cache
//     shared_store: filesystem
//     cache_ttl: 24h

// limits_config:
//   retention_period: ${RETENTION_DAYS}d
//   enforce_metric_name: false
//   reject_old_samples: true
//   reject_old_samples_max_age: 168h
//   max_global_streams_per_user: 5000
//   ingestion_rate_mb: 4
//   ingestion_burst_size_mb: 6

// chunk_store_config:
//   max_look_back_period: 0s

// table_manager:
//   retention_deletes_enabled: true
//   retention_period: ${RETENTION_DAYS}d

// ruler:
//   storage:
//     type: local
//     local:
//       directory: /loki/rules
//   rule_path: /loki/rules
//   alertmanager_url: http://alertmanager:9093
//   ring:
//     kvstore:
//       store: inmemory
//   enable_api: true`

// 	return &corev1.ConfigMap{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      g.options.Name + "-config",
// 			Namespace: g.options.Namespace,
// 			Labels:    g.options.Labels,
// 		},
// 		Data: map[string]string{
// 			"loki.yaml": lokiConfig,
// 		},
// 	}
// }

func (g *ConfigGenerator) GenerateConfigMap() *corev1.ConfigMap {
	lokiConfig := fmt.Sprintf(`
server:
  http_listen_port: 3100
  grpc_listen_port: 9096

auth_enabled: false

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
    cache_ttl: 24h

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    final_sleep: 0s
  chunk_idle_period: 5m
  chunk_retain_period: 30s

limits_config:
  retention_period: %dd
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  max_global_streams_per_user: 5000
  ingestion_rate_mb: 4
  ingestion_burst_size_mb: 6

chunk_store_config:
  max_look_back_period: 0s

table_manager:
  retention_deletes_enabled: true
  retention_period: %dd

ruler:
  storage:
    type: local
    local:
      directory: /loki/rules
  rule_path: /loki/rules
  alertmanager_url: http://alertmanager:9093
  ring:
    kvstore:
      store: inmemory
  enable_api: true`, g.options.RetentionDays, g.options.RetentionDays)

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.options.Name + "-config",
			Namespace: g.options.Namespace,
			Labels:    g.options.Labels,
		},
		Data: map[string]string{
			"loki.yaml": lokiConfig,
		},
	}
}
