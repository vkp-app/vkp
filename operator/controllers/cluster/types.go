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

	EnvStorageClass = "PAAS_STORAGE_CLASS"

	EnvVclusterImage = "RELATED_IMAGE_VCLUSTER_SYNCER"
	EnvCoreDNSImage  = "RELATED_IMAGE_COREDNS"
	EnvSyncImage     = "RELATED_IMAGE_PLUGIN_SYNC"
	EnvHookImage     = "RELATED_IMAGE_PLUGIN_HOOKS"

	EnvPluginPolicy = "PASS_PLUGIN_PULL_POLICY"

	EnvIDPURL = "PAAS_IDP_URL"

	EnvIsOpenShift = "PAAS_IS_OPENSHIFT"
)

type ValuesTemplate struct {
	Name          string
	Ingress       ValuesIngress
	IDP           ValuesIDP
	Storage       paasv1alpha1.Storage
	HA            ValuesHA
	OpenShift     bool
	Image         string
	VclusterImage string
	CoreDNSImage  string
	Plugins       ValuesPlugins
	CustomCA      string
}

type ValuesHA struct {
	Enabled    bool
	Connection string
}

type ValuesPlugins struct {
	SyncImage  string
	HookImage  string
	PullPolicy string
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
