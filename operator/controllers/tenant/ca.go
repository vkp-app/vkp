package tenant

import (
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SecretKeyCA = "ca.crt"

func CASecretName(tenant string) string {
	return fmt.Sprintf("tls-ca-%s", tenant)
}

func CASecret(tenant *paasv1alpha1.Tenant, ca string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CASecretName(tenant.GetName()),
			Namespace: tenant.GetName(),
		},
		Data: map[string][]byte{
			SecretKeyCA: []byte(ca),
		},
	}
}
