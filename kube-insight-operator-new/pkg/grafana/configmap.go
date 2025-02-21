package grafana

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (g *Grafana) GenerateConfigMap() *corev1.ConfigMap {
	datasourcesYaml := `apiVersion: 1
datasources:`

	// Add each data source from the CR
	for _, ds := range g.opts.AdditionalDataSources {
		datasourcesYaml += fmt.Sprintf(`
  - name: %s
    type: %s
    access: proxy
    url: %s
    isDefault: %v
    editable: true`, ds.Name, ds.Type, ds.URL, ds.IsDefault)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.opts.Name + "-config",
			Namespace: g.opts.Namespace,
			Labels:    g.opts.Labels,
		},
		Data: map[string]string{
			"grafana.ini": fmt.Sprintf(`[auth.anonymous]
enabled = false

[security]
admin_user = admin
admin_password = %s`, g.opts.AdminPassword),
			"datasources.yaml": datasourcesYaml,
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
		},
	}
}
