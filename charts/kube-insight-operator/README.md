# Kube Insight Operator

A Helm chart for deploying the Kube Insight Operator, which manages a complete observability stack in Kubernetes.

## Overview

Kube Insight Operator automates the deployment and management of a comprehensive observability solution including:

- **Prometheus**: For metrics collection and monitoring
- **Grafana**: For visualization and dashboards
- **Loki**: For log aggregation
- **Promtail**: For log collection
- **Tempo**: For distributed tracing

This Helm chart simplifies the deployment of the operator, which then manages the entire observability stack lifecycle.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure (for persistent components)

## Installation

### Quick Start

```bash
# Add the repository (if published)
# helm repo add kube-insight-operator https://github.com/johnwroge/kube-insight-operator/releases/
# helm repo update

# Create namespace
kubectl create namespace kube-insight-operator-system

# Install the CRD first
kubectl apply -f https://raw.githubusercontent.com/johnwroge/kube-insight-operator/main/config/crd/bases/monitoring.monitoring.example.com_observabilitystacks.yaml

# Install the chart (from repository)
helm install kube-insight-operator charts/kube-insight-operator \
  --namespace kube-insight-operator-system
```

### From Source

```bash
# Clone the repository
git clone https://github.com/johnwroge/kube-insight-operator.git
cd kube-insight-operator

# Install the CRD
kubectl apply -f config/crd/bases/monitoring.monitoring.example.com_observabilitystacks.yaml

# Install the chart
helm install kube-insight-operator charts/kube-insight-operator \
  --namespace kube-insight-operator-system
```

### Creating an ObservabilityStack

After installing the operator, create an ObservabilityStack custom resource:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: monitoring.monitoring.example.com/v1alpha1
kind: ObservabilityStack
metadata:
  name: monitoring-stack
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
      url: "http://monitoring-stack-prometheus:9090"
      isDefault: true
    - name: "tempo"
      type: "tempo"
      url: "http://monitoring-stack-tempo:3200"
    - name: "loki"
      type: "loki"
      url: "http://monitoring-stack-loki:3100"
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
EOF
```

## Configuration

The following table lists the configurable parameters of the Kube Insight Operator chart and their default values.

### Operator Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `operator.image.repository` | Operator image repository | `johnwroge/kube-insight-operator` |
| `operator.image.tag` | Operator image tag | `latest` |
| `operator.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `operator.resources.requests.cpu` | CPU resource requests | `100m` |
| `operator.resources.requests.memory` | Memory resource requests | `128Mi` |
| `operator.resources.limits.cpu` | CPU resource limits | `500m` |
| `operator.resources.limits.memory` | Memory resource limits | `256Mi` |
| `operator.nodeSelector` | Node labels for pod assignment | `{}` |
| `operator.tolerations` | Node tolerations for pod assignment | `[]` |
| `operator.affinity` | Node affinity for pod assignment | `{}` |
| `operator.securityContext` | Container security context | `{}` |
| `operator.podSecurityContext` | Pod security context | `{}` |

### RBAC Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `rbac.create` | If true, create & use RBAC resources | `true` |
| `serviceAccount.create` | If true, create a service account | `true` |
| `serviceAccount.name` | Name of the service account to use or create | `""` |
| `serviceAccount.annotations` | Annotations for the service account | `{}` |

### Default Stack Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `defaultStack.create` | If true, automatically create a default ObservabilityStack | `false` |
| `defaultStack.spec.*` | Configuration for the default stack | See values.yaml |

## Component Configuration

### Prometheus

The ObservabilityStack allows configuration of Prometheus:

```yaml
prometheus:
  enabled: true
  storage: "10Gi"  # PVC size
  retention: "15d"  # Data retention period
  nodeExporter:
    enabled: true  # Enable node metrics collection
  kubeStateMetrics:
    enabled: true  # Enable Kubernetes state metrics
```

### Grafana

Grafana configuration options:

```yaml
grafana:
  enabled: true
  adminPassword: "admin"  # Default admin password
  serviceType: "ClusterIP"  # Service type (ClusterIP, NodePort, LoadBalancer)
  storage: "5Gi"  # PVC size
  defaultDashboards: true  # Install default dashboards
  additionalDataSources:  # Configure data sources
  - name: "prometheus"
    type: "prometheus"
    url: "http://monitoring-stack-prometheus:9090"
    isDefault: true
```

### Loki

Loki configuration options:

```yaml
loki:
  enabled: true
  storage: "10Gi"  # PVC size
  retentionDays: 15  # Log retention in days
```

### Promtail

Promtail configuration options:

```yaml
promtail:
  enabled: true
  resources:  # Container resources
    cpuRequest: "100m"
    memoryRequest: "128Mi"
    cpuLimit: "200m"
    memoryLimit: "256Mi"
  scrapeKubernetesLogs: true  # Collect Kubernetes logs
```

### Tempo

Tempo configuration options:

```yaml
tempo:
  enabled: true
  storage: "10Gi"  # PVC size
  retentionDays: 7  # Trace retention in days
  resources:  # Container resources
    cpuRequest: "200m"
    memoryRequest: "512Mi"
    cpuLimit: "1"
    memoryLimit: "2Gi"
```

## Accessing Components

After deploying an ObservabilityStack, you can access the components using port-forwarding:

```bash
# Grafana
kubectl port-forward svc/monitoring-stack-grafana 3000:3000
# Access at http://localhost:3000 (default credentials: admin/admin)

# Prometheus
kubectl port-forward svc/monitoring-stack-prometheus 9090:9090
# Access at http://localhost:9090

# Loki
kubectl port-forward svc/monitoring-stack-loki 3100:3100
# Access at http://localhost:3100

# Tempo
kubectl port-forward svc/monitoring-stack-tempo 3200:3200
# Access at http://localhost:3200
```

## Uninstallation

To uninstall the chart:

```bash
helm uninstall kube-insight-operator -n kube-insight-operator-system
```

This will remove the operator but not the CRDs or any deployed observability stacks. To clean up completely:

```bash
# Delete any deployed observability stacks
kubectl delete observabilitystacks --all

# Delete CRDs
kubectl delete crd observabilitystacks.monitoring.monitoring.example.com
```

## Troubleshooting

### Operator Logs

Check the operator logs for issues:

```bash
kubectl logs -n kube-insight-operator-system -l app.kubernetes.io/name=kube-insight-operator
```

### Common Issues

- **PVC Creation Fails**: Ensure your cluster has a default StorageClass
- **RBAC Issues**: The operator needs specific permissions to create and manage resources
- **Image Pull Errors**: Ensure the operator image is accessible to your cluster

## License

This chart is licensed under the Apache License 2.0. See the [LICENSE](https://github.com/johnwroge/kube-insight-operator/blob/main/LICENSE) file for details.