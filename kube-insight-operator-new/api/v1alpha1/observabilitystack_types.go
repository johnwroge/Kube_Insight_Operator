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
type PrometheusSpec struct {
	Enabled   bool   `json:"enabled"`
	Storage   string `json:"storage,omitempty"`
	Retention string `json:"retention,omitempty"`
	// Using named types instead of anonymous structs
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
	// Persistence configuration
	Storage string `json:"storage,omitempty"`
	// Default dashboards to create
	DefaultDashboards bool `json:"defaultDashboards,omitempty"`
	// Additional datasources to configure
	AdditionalDataSources []GrafanaDataSource `json:"additionalDataSources,omitempty"`
}

// GrafanaDataSource defines a data source configuration
type GrafanaDataSource struct {
	// Name of the data source
	Name string `json:"name"`
	// Type of data source (prometheus, loki, tempo, etc)
	Type string `json:"type"`
	// URL of the data source
	URL string `json:"url"`
	// Whether this is the default data source
	IsDefault bool `json:"isDefault,omitempty"`
}

// ObservabilityStackSpec defines the desired state of ObservabilityStack
type ObservabilityStackSpec struct {
	// Prometheus configuration
	Prometheus PrometheusSpec `json:"prometheus,omitempty"`
	// Grafana configuration
	Grafana GrafanaSpec `json:"grafana,omitempty"`
}

// ObservabilityStackStatus defines the observed state of ObservabilityStack
type ObservabilityStackStatus struct {
	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
