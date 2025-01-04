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

## Setup

```bash
# 1. Start Kubernetes (Docker Desktop or minikube)
minikube start  # or start Docker Desktop

# 2. Install CRDs
make install

# 3. Run the operator
make run

# 4. Apply the sample CR
kubectl apply -f config/samples/monitoring_v1alpha1_observabilitystack.yaml
```

# 5. Port-forward all services:
```bash
# In separate terminals:
kubectl port-forward svc/monitoring-test-prometheus 9090:9090
kubectl port-forward svc/monitoring-test-grafana 3000:3000
kubectl port-forward svc/monitoring-test-loki 3100:3100
```

# 6. Access the services:
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (login with admin/admin)
- Loki: http://localhost:3100/ready (should return "ready")