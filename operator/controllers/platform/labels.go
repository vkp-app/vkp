package platform

import paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"

func apiLabels(pr *paasv1alpha1.Platform) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       pr.Name,
		"app.kubernetes.io/instance":   pr.Name,
		"app.kubernetes.io/component":  componentApiServer,
		"app.kubernetes.io/managed-by": "kube-glass-operator",
	}
}

func commonLabels(pr *paasv1alpha1.Platform) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       pr.Name,
		"app.kubernetes.io/instance":   pr.Name,
		"app.kubernetes.io/managed-by": "kube-glass-operator",
	}
}
