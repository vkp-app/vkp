package platform

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"text/template"
)

const (
	portHttp      = "http"
	portGrpc      = "grpc"
	portTelemetry = "telemetry"
)

//go:embed config/dex.tpl.yaml
var dexConfigTemplate string

var dexTpl = template.Must(template.New("config.yaml").Parse(dexConfigTemplate))

type DexConfigTemplate struct {
	Domain string
}

func DexConfig(ctx context.Context, pr *paasv1alpha1.Platform) (*corev1.ConfigMap, error) {
	log := logging.FromContext(ctx)
	// template out the config.yaml
	dexConfig := DexConfigTemplate{
		Domain: pr.Spec.Domain,
	}
	config := new(bytes.Buffer)
	if err := dexTpl.Execute(config, dexConfig); err != nil {
		log.Error(err, "failed to template dex config.yaml")
		return nil, err
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentDex,
			Namespace: pr.Spec.Namespace,
			Labels:    Labels(pr, ComponentDex),
		},
		Data: map[string]string{
			"config.yaml": config.String(),
		},
	}, nil
}

func DexService(pr *paasv1alpha1.Platform) *corev1.Service {
	labels := Labels(pr, ComponentDex)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentDex,
			Namespace: pr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       portHttp,
					Protocol:   corev1.ProtocolTCP,
					Port:       5556,
					TargetPort: intstr.FromString(portHttp),
				},
				{
					Name:       portGrpc,
					Protocol:   corev1.ProtocolTCP,
					Port:       5557,
					TargetPort: intstr.FromString(portGrpc),
				},
				{
					Name:       portTelemetry,
					Protocol:   corev1.ProtocolTCP,
					Port:       5558,
					TargetPort: intstr.FromString(portTelemetry),
				},
			},
			Selector: labels,
		},
	}
}

func DexIngress(pr *paasv1alpha1.Platform) *netv1.Ingress {
	pt := netv1.PathTypePrefix
	className := ingressClassName(&pr.Spec.Dex.Ingress, pr)
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ComponentDex,
			Namespace:   pr.Spec.Namespace,
			Labels:      Labels(pr, ComponentDex),
			Annotations: ingressAnnotations(&pr.Spec.Dex.Ingress, pr),
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &className,
			Rules: []netv1.IngressRule{
				{
					Host: fmt.Sprintf("dex.%s", pr.Spec.Domain),
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pt,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: ComponentDex,
											Port: netv1.ServiceBackendPort{
												Name: portHttp,
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
					SecretName: pr.Spec.Dex.Ingress.SecretRef.Name,
					Hosts: []string{
						fmt.Sprintf("dex.%s", pr.Spec.Domain),
					},
				},
			},
		},
	}
}

func DexDeployment(pr *paasv1alpha1.Platform) *appsv1.Deployment {
	labels := Labels(pr, ComponentDex)
	replicaCount := pr.Spec.Dex.ReplicaCount
	if replicaCount <= 0 {
		replicaCount = 1
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentDex,
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
						annotationMainContainer: ComponentDex,
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets:   pr.Spec.ImagePullSecrets,
					ServiceAccountName: ComponentDex,
					Volumes: []corev1.Volume{
						{
							Name: volumeConfig,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ComponentDex,
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            ComponentDex,
							Image:           pr.Spec.Dex.Image,
							ImagePullPolicy: imagePullPolicy(&pr.Spec.Dex.ComponentSpec, pr),
							Args: append([]string{
								"dex",
								"serve",
								"--web-http-addr",
								"0.0.0.0:5556",
								"--grpc-addr",
								"0.0.0.0:5557",
								"--telemetry-addr",
								"0.0.0.0:5558",
								"/etc/dex/config.yaml",
							}, pr.Spec.Dex.ExtraArgs...),
							Env: []corev1.EnvVar{
								{
									Name: "CLIENT_SECRET_VPK",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: SecretCommonName(pr),
											},
											Key: SecretKeyDexClientSecret,
										},
									},
								},
							},
							Resources: pr.Spec.Dex.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          portHttp,
									ContainerPort: 5556,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          portGrpc,
									ContainerPort: 5557,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          portTelemetry,
									ContainerPort: 5558,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz/ready",
										Port:   intstr.FromString(portTelemetry),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								TimeoutSeconds:   1,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz/live",
										Port:   intstr.FromString(portTelemetry),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								TimeoutSeconds:   1,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeConfig,
									MountPath: "/etc/dex",
									ReadOnly:  true,
								},
							},
						},
					},
				},
			},
		},
	}
}
