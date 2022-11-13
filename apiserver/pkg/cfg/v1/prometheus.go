package v1

type PrometheusMetric struct {
	Name   string `yaml:"name"`
	Metric string `yaml:"metric"`
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
			},
			{
				Name:   "CPU usage",
				Metric: `sum(rate(container_cpu_usage_seconds_total{namespace="%s", pod=~".*-%s|%s-.+"}[1m])) by (namespace)`,
			},
			{
				Name:   "Pod count",
				Metric: `sum by (namespace) (kube_pod_status_ready{namespace="%s", pod=~".*-%s|%s-.+", condition="true"})`,
			},
		},
	}
}
