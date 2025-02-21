# Kube_Insight_Operator

A Kubernetes Operator that automates deployment and management of a complete observability stack, combining metrics, logs, and traces with cost optimization insights.

## Overview

Kube_Insight_Operator simplifies Kubernetes observability by:
- Automating deployment of Prometheus, Loki, Tempo, and Grafana
- Providing pre-configured dashboards and alert rules
- Offering cost optimization recommendations
- Managing the entire observability lifecycle

## Status
This project is currently experimental/under development.

## Features
- One-click observability stack deployment
- Automated configuration and integration
- Pre-built dashboards for common use cases
- Cost analysis and optimization



# Kube Insight Operator

A Kubernetes operator that deploys and manages a complete observability stack including Prometheus, Grafana, Loki, Promtail, and Tempo.

## Overview

This operator provides a unified way to deploy and manage:
- Metrics monitoring with Prometheus
- Log aggregation with Loki and Promtail
- Distributed tracing with Tempo
- Visualization with Grafana

## Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured to communicate with your cluster
- make
- golang 1.19+

## Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/kube-insight-operator
cd kube-insight-operator
```

2. Install the CRDs:
```bash
make manifests
make install
```

3. Run the operator:
```bash
make run
```

## Usage

1. Create an ObservabilityStack by applying the following YAML:

```yaml
apiVersion: monitoring.monitoring.example.com/v1alpha1
kind: ObservabilityStack
metadata:
  name: monitoring-test
spec:
  prometheus:
    enabled: true
    storage: "10Gi"
    retention: "15d"
    nodeExporter:
      enabled: true
    kubeStateMetrics:
      enabled: true
  grafana:
    enabled: true
    adminPassword: "admin"
    serviceType: "ClusterIP"
    storage: "5Gi"
    defaultDashboards: true
    additionalDataSources:
    - name: "prometheus"
      type: "prometheus"
      url: "http://monitoring-test-prometheus:9090"
      isDefault: true
    - name: "tempo"
      type: "tempo"
      url: "http://monitoring-test-tempo:3200"
    - name: "loki"
      type: "loki"
      url: "http://monitoring-test-loki:3100"
  loki:
    enabled: true
    storage: "10Gi"
    retentionDays: 15
  promtail:
    enabled: true
    resources:
      cpuRequest: "100m"
      memoryRequest: "128Mi"
      cpuLimit: "200m"
      memoryLimit: "256Mi"
    scrapeKubernetesLogs: true
  tempo:
    enabled: true
    storage: "10Gi"
    retentionDays: 7
    resources:
      cpuRequest: "200m"
      memoryRequest: "512Mi"
      cpuLimit: "1"
      memoryLimit: "2Gi"
```

2. Apply the configuration:
```bash
kubectl apply -f config/samples/monitoring_v1alpha1_observabilitystack.yaml
```

3. Access the components:

```bash
# Grafana
kubectl port-forward svc/monitoring-test-grafana 3000:3000
# Access at http://localhost:3000 (default credentials: admin/admin)

# Prometheus
kubectl port-forward svc/monitoring-test-prometheus 9090:9090
# Access at http://localhost:9090

# Loki
kubectl port-forward svc/monitoring-test-loki 3100:3100
# Access at http://localhost:3100

# Tempo
kubectl port-forward svc/monitoring-test-tempo 3200:3200
# Access at http://localhost:3200
```

## Configuration Options

### Prometheus
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Prometheus | false |
| storage | Storage size | "10Gi" |
| retention | Data retention period | "15d" |
| nodeExporter.enabled | Enable node exporter | true |
| kubeStateMetrics.enabled | Enable kube-state-metrics | true |

### Grafana
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Grafana | false |
| adminPassword | Admin password | "admin" |
| serviceType | Service type | "ClusterIP" |
| storage | Storage size | "5Gi" |
| defaultDashboards | Enable default dashboards | true |
| additionalDataSources | Additional data sources | [] |

### Loki
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Loki | false |
| storage | Storage size | "10Gi" |
| retentionDays | Log retention period in days | 14 |

### Promtail
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Promtail | false |
| resources | Resource requests and limits | see example |
| scrapeKubernetesLogs | Enable Kubernetes log scraping | true |

### Tempo
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Tempo | false |
| storage | Storage size | "10Gi" |
| retentionDays | Trace retention period in days | 7 |
| resources | Resource requests and limits | see example |

## Development

1. Make changes to the operator code
2. Update CRDs:
```bash
make manifests
```
3. Install updated CRDs:
```bash
make install
```
4. Run the operator locally:
```bash
make run
```

## Troubleshooting

Common issues and solutions:

1. Pods not starting:
```bash
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

2. PVC issues:
```bash
kubectl get pvc
kubectl describe pvc <pvc-name>
```

3. Check operator logs:
```bash
kubectl logs -l app=kube-insight-operator
```


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

