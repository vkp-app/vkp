package v1

import "gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"

type PrometheusMetric struct {
	Name   string             `yaml:"name"`
	Metric string             `yaml:"metric"`
	Format model.MetricFormat `yaml:"format"`
}

type PrometheusConfig struct {
	ClusterMetrics []PrometheusMetric `yaml:"clusterMetrics"`
}

// NewPrometheusConfig provides default configuration
// that is used if the user doesn't override anything
func NewPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		ClusterMetrics: []PrometheusMetric{
			{
				Name:   "Memory usage",
				Metric: `sum by (namespace) (topk(1, container_memory_usage_bytes{namespace="{namespace}"}) * on (pod,namespace) group_right kube_pod_labels{label_paas_dcas_dev_metric_target="{cluster}"})`,
				Format: model.MetricFormatBytes,
			},
			{
				Name:   "CPU usage",
				Metric: `sum(rate(container_cpu_usage_seconds_total{namespace="{namespace}", pod=~".*-{cluster}|{cluster}-.+"}[1m])) by (namespace)`,
				Format: model.MetricFormatCPU,
			},
			{
				Name:   "Pod count",
				Metric: `sum by (namespace) (kube_pod_status_ready{namespace="{namespace}", condition="true"} * on (pod,namespace) group_right kube_pod_labels{label_paas_dcas_dev_metric_target="{cluster}"})`,
				Format: model.MetricFormatPlain,
			},
			{
				Name:   "Request volume",
				Metric: `sum by (exported_namespace) (irate(nginx_ingress_controller_requests{exported_namespace="%s", exported_service=~".*-%s|%s-.+"}[2m]))`,
				Format: model.MetricFormatRps,
			},
		},
	}
}
