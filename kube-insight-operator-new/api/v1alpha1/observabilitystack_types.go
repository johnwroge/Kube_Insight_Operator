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
type PrometheusSpec struct {
	// Whether Prometheus is enabled
	Enabled bool `json:"enabled"`
	// Storage size (e.g., "10Gi")
	Storage string `json:"storage,omitempty"`
	// Retention period (e.g., "15d")
	Retention string `json:"retention,omitempty"`
}

// GrafanaSpec defines the configuration for Grafana
type GrafanaSpec struct {
	// Whether Grafana is enabled
	Enabled bool `json:"enabled"`
	// Admin password for Grafana
	AdminPassword string `json:"adminPassword,omitempty"`
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
