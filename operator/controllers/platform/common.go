package platform

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SecretCommonName(pr *paasv1alpha1.Platform) string {
	return fmt.Sprintf("%s-common", pr.GetName())
}

func CommonSecrets(pr *paasv1alpha1.Platform) *corev1.Secret {
	// generate a cookie string
	// https://stackoverflow.com/a/65607935
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	cookieSecret := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%x", b)[:32]))

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretCommonName(pr),
			Namespace: pr.Spec.Namespace,
			Labels:    commonLabels(pr),
		},
		Data: map[string][]byte{
			SecretKeyOauthCookie: []byte(cookieSecret),
		},
	}
}

func imagePullPolicy(c *paasv1alpha1.ComponentSpec, pr *paasv1alpha1.Platform) corev1.PullPolicy {
	if c.ImagePullPolicy != "" {
		return c.ImagePullPolicy
	}
	return pr.Spec.ImagePullPolicy
}
