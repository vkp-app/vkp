package cluster

import (
	"fmt"
	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func IdentitySecretName(cluster string) string {
	return fmt.Sprintf("tls-identity-%s", cluster)
}

func Certificate(cr *paasv1alpha1.Cluster, issuer, issuerKind string) *certv1.Certificate {
	return &certv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IdentitySecretName(cr.GetName()),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
		},
		Spec: certv1.CertificateSpec{
			CommonName:  fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()),
			Duration:    &metav1.Duration{Duration: time.Hour * 8760},
			RenewBefore: &metav1.Duration{Duration: time.Hour * 2190},
			IsCA:        false,
			PrivateKey: &certv1.CertificatePrivateKey{
				Algorithm: certv1.ECDSAKeyAlgorithm,
			},
			SecretName: IdentitySecretName(cr.GetName()),
			IssuerRef: certmetav1.ObjectReference{
				Name:  issuer,
				Kind:  issuerKind,
				Group: "cert-manager.io",
			},
			Usages: []certv1.KeyUsage{
				certv1.UsageDigitalSignature,
				certv1.UsageKeyEncipherment,
				certv1.UsageClientAuth,
			},
		},
	}
}
