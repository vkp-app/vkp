package cluster

import paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"

const (
	EnvHostname     = "PAAS_HOSTNAME"
	EnvChartName    = "PAAS_CHART_NAME"
	EnvChartRepo    = "PAAS_CHART_REPO"
	EnvChartVersion = "PAAS_CHART_VERSION"

	EnvKubeVersion = "PAAS_KUBE_VERSION"

	EnvIngressClass  = "PAAS_INGRESS_CLASS"
	EnvIngressIssuer = "PAAS_INGRESS_ISSUER"

	EnvIDPURL = "PAAS_IDP_URL"

	EnvIsOpenShift = "PAAS_IS_OPENSHIFT"
)

type ValuesTemplate struct {
	Name      string
	Ingress   ValuesIngress
	IDP       ValuesIDP
	Storage   paasv1alpha1.Storage
	HA        bool
	OpenShift bool
	Image     string
}

type ValuesIngress struct {
	Host          string
	TLSSecretName string
	ClassName     string
}

type ValuesIDP struct {
	URL        string
	SecretName string
	CustomCA   bool
}
