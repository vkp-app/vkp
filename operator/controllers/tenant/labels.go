package tenant

import (
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
)

func Labels(tr *v1alpha1.Tenant) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       tr.Name,
		"app.kubernetes.io/instance":   tr.Name,
		"app.kubernetes.io/component":  tr.Kind,
		"app.kubernetes.io/managed-by": "vkp",
		LabelOwned:                     "true",
	}
}
