package cluster

import paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"

func Labels(cr *paasv1alpha1.Cluster) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/component":  cr.Kind,
		"app.kubernetes.io/managed-by": "kube-glass-operator",
	}
}