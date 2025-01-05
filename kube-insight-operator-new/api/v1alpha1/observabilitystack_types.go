/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrometheusSpec defines the configuration for Prometheus
// type PrometheusSpec struct {
//     Enabled bool `json:"enabled"`
//     Storage string `json:"storage,omitempty"`
//     Retention string `json:"retention,omitempty"`
//     // Adding node exporter configuration
//     NodeExporter struct {
//         Enabled bool `json:"enabled"`
//     } `json:"nodeExporter,omitempty"`
//     // Adding kube-state-metrics configuration
//     KubeStateMetrics struct {
//         Enabled bool `json:"enabled"`
//     } `json:"kubeStateMetrics,omitempty"`
// }

// NodeExporterSpec defines the configuration for node-exporter
type NodeExporterSpec struct {
	Enabled bool `json:"enabled"`
}

// KubeStateMetricsSpec defines the configuration for kube-state-metrics
type KubeStateMetricsSpec struct {
	Enabled bool `json:"enabled"`
}

// PrometheusSpec defines the configuration for Prometheus
//
//	type PrometheusSpec struct {
//		Enabled   bool   `json:"enabled"`
//		Storage   string `json:"storage,omitempty"`
//		Retention string `json:"retention,omitempty"`
//		// Using named types instead of anonymous structs
//		NodeExporter     NodeExporterSpec     `json:"nodeExporter,omitempty"`
//		KubeStateMetrics KubeStateMetricsSpec `json:"kubeStateMetrics,omitempty"`
//	}
type PrometheusSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[0-9]+[GM]i$`
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[0-9]+[hdw]$`
	Retention string `json:"retention,omitempty"`

	NodeExporter     NodeExporterSpec     `json:"nodeExporter,omitempty"`
	KubeStateMetrics KubeStateMetricsSpec `json:"kubeStateMetrics,omitempty"`
}

// GrafanaSpec defines the configuration for Grafana
type GrafanaSpec struct {
	// Whether Grafana is enabled
	Enabled bool `json:"enabled"`
	// Admin password for Grafana
	AdminPassword string `json:"adminPassword,omitempty"`
	// Service type (LoadBalancer, ClusterIP, NodePort)
	ServiceType string `json:"serviceType,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[0-9]+[GM]i$`
	Storage string `json:"storage,omitempty"`
	// Default dashboards to create
	DefaultDashboards bool `json:"defaultDashboards,omitempty"`
	// Additional datasources to configure
	AdditionalDataSources []GrafanaDataSource `json:"additionalDataSources,omitempty"`
}

// GrafanaDataSource defines a data source configuration
type GrafanaDataSource struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=prometheus;loki;tempo
	Type string `json:"type"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https?://.*`
	URL string `json:"url"`

	// +kubebuilder:validation:Optional
	IsDefault bool `json:"isDefault,omitempty"`
}

type ResourceRequirements struct {
    // +kubebuilder:validation:Optional
    // +kubebuilder:default="100m"
    CPURequest string `json:"cpuRequest,omitempty"`

    // +kubebuilder:validation:Optional
    // +kubebuilder:default="128Mi"
    MemoryRequest string `json:"memoryRequest,omitempty"`

    // +kubebuilder:validation:Optional
    // +kubebuilder:default="200m"
    CPULimit string `json:"cpuLimit,omitempty"`

    // +kubebuilder:validation:Optional
    // +kubebuilder:default="256Mi"
    MemoryLimit string `json:"memoryLimit,omitempty"`
}

type PromtailSpec struct {
    // +kubebuilder:validation:Optional
    // +kubebuilder:default=false
    Enabled bool `json:"enabled"`

    // +kubebuilder:validation:Optional
    Resources ResourceRequirements `json:"resources,omitempty"`

    // +kubebuilder:validation:Optional
    Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

    // +kubebuilder:validation:Optional
    // +kubebuilder:default=true
    ScrapeKubernetesLogs bool `json:"scrapeKubernetesLogs,omitempty"`

    // +kubebuilder:validation:Optional
    ExtraArgs []string `json:"extraArgs,omitempty"`
}

// ObservabilityStackSpec defines the desired state of ObservabilityStack
type ObservabilityStackSpec struct {
	// Prometheus configuration
	Prometheus PrometheusSpec `json:"prometheus,omitempty"`
	// Grafana configuration
	Grafana GrafanaSpec `json:"grafana,omitempty"`

	// +kubebuilder:validation:Optional
	Loki LokiSpec `json:"loki,omitempty"`
	
	Promtail PromtailSpec `json:"promtail,omitempty"`

}

// ObservabilityStackStatus defines the observed state of ObservabilityStack
type ObservabilityStackStatus struct {
	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type LokiSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="10Gi"
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=14
	// +kubebuilder:validation:Minimum=1
	RetentionDays int32 `json:"retentionDays,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ObservabilityStack is the Schema for the observabilitystacks API
type ObservabilityStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObservabilityStackSpec   `json:"spec,omitempty"`
	Status ObservabilityStackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ObservabilityStackList contains a list of ObservabilityStack
type ObservabilityStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObservabilityStack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObservabilityStack{}, &ObservabilityStackList{})
}
