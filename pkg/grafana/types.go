package grafana

import (
	monitoringv1alpha1 "github.com/johnwroge/kube-insight-operator/kube-insight-operator-new/api/v1alpha1"
)

type Options struct {
	Name                  string
	Namespace             string
	Labels                map[string]string
	AdminPassword         string
	Storage               string
	AdditionalDataSources []monitoringv1alpha1.GrafanaDataSource
	DefaultDashboards     bool
}

type Grafana struct {
	opts *Options
}

func New(opts Options) *Grafana {
	return &Grafana{
		opts: &opts,
	}
}
