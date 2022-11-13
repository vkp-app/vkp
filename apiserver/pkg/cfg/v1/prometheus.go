package v1

import "gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"

type PrometheusMetric struct {
	Name   string             `yaml:"name"`
	Metric string             `yaml:"metric"`
	Format model.MetricFormat `yaml:"format"`
}

type PrometheusConfig struct {
	Metrics []PrometheusMetric `yaml:"metrics"`
}

// NewPrometheusConfig provides default configuration
// that is used if the user doesn't override anything
func NewPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Metrics: []PrometheusMetric{
			{
				Name:   "Memory usage",
				Metric: `sum by (namespace) (container_memory_usage_bytes{namespace="%s", pod=~".*-%s|%s-.+"})`,
				Format: model.MetricFormatBytes,
			},
			{
				Name:   "CPU usage",
				Metric: `sum(rate(container_cpu_usage_seconds_total{namespace="%s", pod=~".*-%s|%s-.+"}[1m])) by (namespace)`,
				Format: model.MetricFormatCPU,
			},
			{
				Name:   "Pod count",
				Metric: `sum by (namespace) (kube_pod_status_ready{namespace="%s", pod=~".*-%s|%s-.+", condition="true"})`,
				Format: model.MetricFormatPlain,
			},
		},
	}
}
