package syncers

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func init() {
	_ = paasv1alpha1.AddToScheme(plugin.Scheme)
	_ = paasv1alpha1.AddToScheme(clientgoscheme.Scheme)

	_ = idpv1.AddToScheme(plugin.Scheme)
	_ = idpv1.AddToScheme(clientgoscheme.Scheme)
}
