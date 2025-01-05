package tempo

import (
	corev1 "k8s.io/api/core/v1"
)

type Options struct {
	Name          string
	Namespace     string
	Labels        map[string]string
	Storage       string
	RetentionDays int32
	Resources     *corev1.ResourceRequirements
}

type ConfigGenerator struct {
	options Options
}

func NewConfigGenerator(opts Options) *ConfigGenerator {
	return &ConfigGenerator{
		options: opts,
	}
}
