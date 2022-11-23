package tenant

import (
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Namespace(tr *paasv1alpha1.Tenant) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: tr.GetName(),
		},
	}
}
