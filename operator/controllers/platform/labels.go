package platform

import paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"

// Deprecated
func apiLabels(pr *paasv1alpha1.Platform) map[string]string {
	return Labels(pr, ComponentApiServer)
}

func Labels(pr *paasv1alpha1.Platform, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       pr.Name,
		"app.kubernetes.io/instance":   pr.Name,
		"app.kubernetes.io/component":  component,
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
