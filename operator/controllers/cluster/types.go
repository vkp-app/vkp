package cluster

const (
	EnvHostname     = "PAAS_HOSTNAME"
	EnvChartName    = "PAAS_CHART_NAME"
	EnvChartRepo    = "PAAS_CHART_REPO"
	EnvChartVersion = "PAAS_CHART_VERSION"

	EnvKubeVersion = "PAAS_KUBE_VERSION"

	EnvIngressClass  = "PAAS_INGRESS_CLASS"
	EnvIngressIssuer = "PAAS_INGRESS_ISSUER"
)

type ValuesTemplate struct {
	Ingress ValuesIngress
}

type ValuesIngress struct {
	ClassName string
	Issuer    string
	Host      string
}
