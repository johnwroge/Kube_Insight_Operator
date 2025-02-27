# Kube Insight Operator Helm Chart

This Helm chart deploys the Kube Insight Operator, which manages observability stacks including Prometheus, Grafana, Loki, Promtail, and Tempo in Kubernetes clusters.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+

## Installing the Chart

First, clone the repository:

```bash
git clone https://github.com/johnwroge/kube-insight-operator.git
cd kube-insight-operator
```

To install the chart with the release name `kube-insight-operator`:

```bash
# Create namespace
kubectl create namespace kube-insight-operator-system

# Install the chart
helm install kube-insight-operator ./charts/kube-insight-operator \
  --namespace kube-