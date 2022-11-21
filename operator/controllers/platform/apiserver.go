package platform

import (
	_ "embed"
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

const (
	volumeConfig = "config"
	volumeDexTLS = "dex-tls"

	portPublic = "public"
)

//go:embed config/prometheus.default.yaml
var defaultPrometheusYAML string

func ApiConfig(pr *paasv1alpha1.Platform) *corev1.ConfigMap {
	prom := pr.Spec.ApiServer.PrometheusConfig
	if prom == "" {
		prom = defaultPrometheusYAML
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentApiServer,
			Namespace: pr.Spec.Namespace,
			Labels:    apiLabels(pr),
		},
		Data: map[string]string{
			"prometheus.yaml": prom,
		},
	}
}

func ApiService(pr *paasv1alpha1.Platform) *corev1.Service {
	labels := apiLabels(pr)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentApiServer,
			Namespace: pr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       portPublic,
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromString(portPublic),
				},
			},
			Selector: labels,
		},
	}
}

func ApiIngress(pr *paasv1alpha1.Platform) *netv1.Ingress {
	pt := netv1.PathTypePrefix
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ComponentApiServer,
			Namespace:   pr.Spec.Namespace,
			Labels:      Labels(pr, ComponentApiServer),
			Annotations: pr.Spec.Ingress.Annotations,
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &pr.Spec.Ingress.IngressClassName,
			Rules: []netv1.IngressRule{
				{
					Host: fmt.Sprintf("vkp.%s", pr.Spec.Domain),
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/api",
									PathType: &pt,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: ComponentApiServer,
											Port: netv1.ServiceBackendPort{
												Name: portPublic,
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
					SecretName: pr.Spec.Ingress.SecretRef.Name,
					Hosts: []string{
						fmt.Sprintf("vkp.%s", pr.Spec.Domain),
					},
				},
			},
		},
	}
}

func ApiDeployment(pr *paasv1alpha1.Platform) *appsv1.Deployment {
	labels := Labels(pr, ComponentApiServer)
	replicaCount := pr.Spec.ApiServer.ReplicaCount
	if replicaCount <= 0 {
		replicaCount = 1
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentApiServer,
			Namespace: pr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(replicaCount),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						annotationMainContainer: ComponentApiServer,
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets:   pr.Spec.ImagePullSecrets,
					ServiceAccountName: ComponentApiServer,
					Volumes: []corev1.Volume{
						{
							Name: volumeConfig,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ComponentApiServer,
									},
								},
							},
						},
						{
							Name: volumeDexTLS,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: pr.Spec.Dex.Ingress.SecretRef.Name,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            ComponentApiServer,
							Image:           pr.Spec.ApiServer.Image,
							ImagePullPolicy: imagePullPolicy(&pr.Spec.ApiServer.ComponentSpec, pr),
							Args: append([]string{
								"--v=5",
								fmt.Sprintf("--prometheus-url=%s", pr.Spec.PrometheusURL),
								"--prometheus-config-file=/etc/vkp/config/prometheus.yaml",
								fmt.Sprintf("--dex-url=https://dex.%s", pr.Spec.Domain),
							}, pr.Spec.ApiServer.ExtraArgs...),
							Env: []corev1.EnvVar{
								{
									Name: "KUBERNETES_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
							},
							Resources: pr.Spec.ApiServer.Resources,
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/readyz",
										Port:   intstr.FromInt(8081),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      3,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/livez",
										Port:   intstr.FromInt(8081),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 15,
								TimeoutSeconds:      15,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeConfig,
									MountPath: "/etc/vkp/config",
									ReadOnly:  true,
								},
							},
						},
						{
							Name:            ComponentOauthProxy,
							Image:           pr.Spec.ApiServer.OauthProxy.Image,
							ImagePullPolicy: imagePullPolicy(&pr.Spec.ApiServer.OauthProxy, pr),
							Args: append([]string{
								"--http-address=:8079",
								"--provider=oidc",
								fmt.Sprintf("--client-id=%s", "todo"),
								fmt.Sprintf("--client-secret=%s", "todo"),
								"--email-domain=*",
								fmt.Sprintf("--oidc-issuer-url=https://dex.%s", pr.Spec.Domain),
								fmt.Sprintf("--redirect-url=https://vkp.%s/oauth2/callback", pr.Spec.Domain),
								"--prefer-email-to-user=true",
								"--cookie-secure=true",
								"--cookie-refresh=2h",
								"--code-challenge-method=S256",
								"--upstream=http://localhost:8080",
								"--scope=openid profile email groups",
							}, pr.Spec.ApiServer.OauthProxy.ExtraArgs...),
							Env: []corev1.EnvVar{
								{
									Name: "OAUTH2_PROXY_COOKIE_SECRET",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: SecretCommonName(pr),
											},
											Key: SecretKeyOauthCookie,
										},
									},
								},
								{
									Name:  "SSL_CERT_DIR",
									Value: "/var/run/secrets/paas.dcas.dev/ca",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          portPublic,
									ContainerPort: 8079,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeDexTLS,
									MountPath: "/var/run/secrets/paas.dcas.dev/ca/ca.crt",
									ReadOnly:  true,
									SubPath:   "ca.crt",
								},
							},
							Resources: pr.Spec.ApiServer.OauthProxy.Resources,
						},
					},
				},
			},
		},
	}
}
