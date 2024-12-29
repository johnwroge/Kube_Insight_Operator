package grafana

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (g *Grafana) GenerateConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.opts.Name + "-config",
			Namespace: g.opts.Namespace,
			Labels:    g.opts.Labels,
		},
		Data: map[string]string{
			"grafana.ini": `
[auth.anonymous]
enabled = false

[security]
admin_user = admin
admin_password = ` + g.opts.AdminPassword + `
`,
			"datasources.yaml": `
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: ` + g.opts.PrometheusURL + `
    isDefault: true
    editable: true
`,
		},
	}
}

func (g *Grafana) GenerateDefaultDashboardsConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.opts.Name + "-dashboards",
			Namespace: g.opts.Namespace,
			Labels:    g.opts.Labels,
		},
		Data: map[string]string{
			"kubernetes-cluster-monitoring.json": `{
                "dashboard": {
                    "title": "Kubernetes Cluster Monitoring",
                    "uid": "kubernetes-cluster",
                    "panels": [
                        // Dashboard JSON will go here
                    ]
                }
            }`,
			// Add more default dashboards as needed
		},
	}
}
