package cluster

import (
	"fmt"
	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	DexKeyID     = "client_id"
	DexKeySecret = "client_secret"
	DexKeyCA     = "ca.crt"
)

func DexSecretName(cluster string) string {
	return fmt.Sprintf("%s-dex", cluster)
}

func DexSecret(cr *paasv1alpha1.Cluster, ca string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DexSecretName(cr.GetName()),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
		},
		Data: map[string][]byte{
			DexKeyID:     []byte(cr.GetUID()),
			DexKeySecret: []byte(uuid.NewUUID()),
			DexKeyCA:     []byte(ca),
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func DexOAuthClient(cr *paasv1alpha1.Cluster, sec *corev1.Secret) *idpv1.OAuthClient {
	return &idpv1.OAuthClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
		},
		Spec: idpv1.OAuthClientSpec{
			ClientID: string(sec.Data[DexKeyID]),
			ClientSecretRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: sec.GetName(),
				},
				Key: DexKeySecret,
			},
			RedirectURIs: []string{
				fmt.Sprintf("https://console.%s.%s/auth/callback", cr.Status.ClusterID, cr.Status.ClusterDomain),
				fmt.Sprintf("https://console.%s.%s/oauth2/callback", cr.Status.ClusterID, cr.Status.ClusterDomain),
				// needed for kubectl oidc-login plugin
				"http://localhost:8000",
			},
		},
	}
}
