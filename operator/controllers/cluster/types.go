package cluster

const (
	EnvHostname     = "PAAS_HOSTNAME"
	EnvChartName    = "PAAS_CHART_NAME"
	EnvChartRepo    = "PAAS_CHART_REPO"
	EnvChartVersion = "PAAS_CHART_VERSION"

	EnvKubeVersion = "PAAS_KUBE_VERSION"

	EnvIngressClass  = "PAAS_INGRESS_CLASS"
	EnvIngressIssuer = "PAAS_INGRESS_ISSUER"

	EnvIDPURL      = "PAAS_IDP_URL"
	EnvIDPClientID = "PAAS_IDP_CLIENT_ID"
)

type ValuesTemplate struct {
	Ingress ValuesIngress
	IDP     ValuesIDP
}

type ValuesIngress struct {
	ClassName     string
	Issuer        string
	Host          string
	TLSSecretName string
}

type ValuesIDP struct {
	URL      string
	ClientID string
}
