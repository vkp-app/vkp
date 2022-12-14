package cluster

import (
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func getHostname(cr *paasv1alpha1.Cluster) string {
	return fmt.Sprintf("api.%s.%s", cr.Status.ClusterID, cr.Status.ClusterDomain)
}

func IngressSecretName(cluster string) string {
	return fmt.Sprintf("tls-kubeapi-%s", cluster)
}

func Ingress(cr *paasv1alpha1.Cluster) *netv1.Ingress {
	hostname := getHostname(cr)
	pathType := netv1.PathTypePrefix
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
			Annotations: map[string]string{
				"cert-manager.io/cluster-issuer":               GetEnv(EnvIngressIssuer, ""),
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/ssl-redirect":     "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: pointer.String(GetEnv(EnvIngressClass, "nginx")),
			Rules: []netv1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: cr.GetName(),
											Port: netv1.ServiceBackendPort{
												Name: "https",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []netv1.IngressTLS{
				{
					SecretName: IngressSecretName(cr.GetName()),
					Hosts: []string{
						hostname,
					},
				},
			},
		},
	}
}
