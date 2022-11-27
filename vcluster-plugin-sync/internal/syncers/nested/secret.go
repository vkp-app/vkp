package nested

import (
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/tenant"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const caSecretName = "platform-root-ca.crt"
const caSecretDescription = "Contains a CA bundle that can be used to verify VKP components and other internal components, depending on how your network was setup."

func RootCA(namespace, ca string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        caSecretName,
			Namespace:   namespace,
			Annotations: map[string]string{"kubernetes.io/description": caSecretDescription},
		},
		Data: map[string]string{
			tenant.SecretKeyCA: ca,
		},
	}
}
