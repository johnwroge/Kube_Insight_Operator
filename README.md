# Kube Insight Operator

A Kubernetes operator that deploys and manages a complete observability stack combining metrics, logs, and traces including Prometheus, Grafana, Loki, Promtail, and Tempo.

## Overview

Kube Insight Operator simplifies Kubernetes observability by:
- Automating deployment of Prometheus, Loki, Tempo, and Grafana
- Providing pre-configured dashboards and alert rules
- Managing the entire observability lifecycle

## Status

**This project is currently experimental/under development.**

### Immediate Priorities
- [ ] Implement comprehensive status reporting for all observability components
- [ ] Develop Stronger Authentication and Authorization Mechanisms
- [ ] Enhance error handling and logging
- [ ] Develop a Helm chart for easier deployment
- [ ] Create detailed documentation for configuration and customization
- [ ] Implement more robust testing suite

## Technology Stack

This project leverages the Kubernetes Operator SDK, a toolkit for developing Kubernetes operators that provides:
- Automated scaffolding and code generation for Kubernetes operator development
- Controller runtime for managing and reconciling custom resources
- Custom Resource Definition (CRD) management
- Reconciliation loop implementation
- Built on the Kubebuilder framework
- Built with Go, the primary language for Kubernetes operator development


### Observability Stack Components
The operator integrates a comprehensive observability solution with:
- **Prometheus**: Comprehensive metrics collection and monitoring system for capturing performance and operational data
- **Grafana**: Advanced visualization platform for creating interactive dashboards and graphical representations of complex metrics
- **Loki**: Scalable log aggregation and storage solution designed for cloud-native environments
- **Promtail**: Efficient log collection agent for gathering logs from various sources
- **Tempo**: Distributed tracing system for tracking and analyzing request flows across microservices

### Core Technologies
- Kubernetes Operator SDK
- Go programming language
- Kubebuilder framework
- Custom Resource Definitions (CRDs)


## Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured to communicate with your cluster
- make
- golang 1.19+

## Installation

1. Clone the repository:
```bash
git clone https://github.com/johnwroge/kube-insight-operator
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
| enabled | Enable Prometheus | true |
| storage | Storage size | "10Gi" |
| retention | Data retention period | "15d" |
| nodeExporter.enabled | Enable node exporter | true |
| kubeStateMetrics.enabled | Enable kube-state-metrics | true |

### Grafana
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Grafana | true |
| adminPassword | Admin password | "admin" |
| serviceType | Service type | "ClusterIP" |
| storage | Storage size | "5Gi" |
| defaultDashboards | Enable default dashboards | true |
| additionalDataSources | Additional data sources | [] |

### Loki
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Loki | true |
| storage | Storage size | "10Gi" |
| retentionDays | Log retention period in days | 14 |

### Promtail
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Promtail | true |
| resources | Resource requests and limits | see example |
| scrapeKubernetesLogs | Enable Kubernetes log scraping | true |

### Tempo
| Parameter | Description | Default |
|-----------|-------------|---------|
| enabled | Enable Tempo | true |
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
This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
