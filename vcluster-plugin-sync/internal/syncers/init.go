package syncers

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func init() {
	_ = paasv1alpha1.AddToScheme(plugin.Scheme)
	_ = paasv1alpha1.AddToScheme(clientgoscheme.Scheme)
}
