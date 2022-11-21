package platform

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// randStr generates a random string
// of a given length.
//
// // https://stackoverflow.com/a/65607935
func randStr(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%x", b)[:length]))
}

func SecretCommonName(pr *paasv1alpha1.Platform) string {
	return fmt.Sprintf("%s-common", pr.GetName())
}

func CommonSecrets(pr *paasv1alpha1.Platform) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretCommonName(pr),
			Namespace: pr.Spec.Namespace,
			Labels:    commonLabels(pr),
		},
		Data: map[string][]byte{
			SecretKeyOauthCookie:     []byte(randStr(32)),
			SecretKeyDexClientSecret: []byte(randStr(32)),
		},
	}
}

func ServiceAccount(c string, pr *paasv1alpha1.Platform) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c,
			Namespace: pr.Spec.Namespace,
			Labels:    Labels(pr, c),
		},
	}
}

func imagePullPolicy(c *paasv1alpha1.ComponentSpec, pr *paasv1alpha1.Platform) corev1.PullPolicy {
	if c.ImagePullPolicy != "" {
		return c.ImagePullPolicy
	}
	return pr.Spec.ImagePullPolicy
}

func ingressClassName(c *paasv1alpha1.ComponentIngressSpec, pr *paasv1alpha1.Platform) string {
	if c.IngressClassName != "" {
		return c.IngressClassName
	}
	return pr.Spec.Ingress.IngressClassName
}

func ingressAnnotations(c *paasv1alpha1.ComponentIngressSpec, pr *paasv1alpha1.Platform) map[string]string {
	if c.Annotations != nil && len(c.Annotations) > 0 {
		return c.Annotations
	}
	return pr.Spec.Ingress.Annotations
}
