package syncers

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

func init() {
	_ = paasv1alpha1.AddToScheme(plugin.Scheme)
	_ = paasv1alpha1.AddToScheme(clientgoscheme.Scheme)

	_ = idpv1.AddToScheme(plugin.Scheme)
	_ = idpv1.AddToScheme(clientgoscheme.Scheme)

	_ = apiextv1.AddToScheme(clientgoscheme.Scheme)
	_ = apiextv1.AddToScheme(plugin.Scheme)

	_ = apiregv1.AddToScheme(clientgoscheme.Scheme)
	_ = apiregv1.AddToScheme(plugin.Scheme)
}
