package cluster

import (
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	DexKeyID     = "client_id"
	DexKeySecret = "client_secret"
)

func DexSecret(cr *paasv1alpha1.Cluster) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-dex", cr.GetName()),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
		},
		Data: map[string][]byte{
			DexKeyID:     []byte(cr.GetUID()),
			DexKeySecret: []byte(uuid.NewUUID()),
		},
		Type: corev1.SecretTypeOpaque,
	}
}
