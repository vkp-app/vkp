package syncers

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

func init() {
	_ = paasv1alpha1.AddToScheme(plugin.Scheme)
}
