package tenant

import (
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	imageDashboardKube = "ghcr.io/vkp-app/addons/dashboard-k8s:1.2.0"
	imageDashboardOKD  = "ghcr.io/vkp-app/addons/dashboard-okd:1.2.0"
	imagePodInfo       = "ghcr.io/vkp-app/addons/podinfo:1.0.1"
	imagePromAdapter   = "ghcr.io/vkp-app/addons/prometheus-adapter:1.4.0"
)

func Addons(tr *v1alpha1.Tenant) []v1alpha1.ClusterAddon {
	labels := Labels(tr)
	return []v1alpha1.ClusterAddon{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dashboard-k8s",
				Namespace: tr.GetName(),
				Labels:    labels,
			},
			Spec: v1alpha1.ClusterAddonSpec{
				Resources: []v1alpha1.RemoteRef{
					{
						OCI: v1alpha1.OCIRemoteRef{
							Name: cluster.GetEnv(cluster.EnvAddonDashboardKubeImage, imageDashboardKube),
						},
					},
				},
				DisplayName: "Kubernetes Dashboard",
				Maintainer:  "The Kubernetes Authors",
				Logo:        "https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo.svg",
				Description: "General-purpose web UI for Kubernetes clusters (Mutually-exclusive with the OpenShift Console).",
				Source:      v1alpha1.SourceCommunity,
				SourceURL:   "https://github.com/kubernetes/dashboard",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dashboard-okd",
				Namespace: tr.GetName(),
				Labels:    labels,
			},
			Spec: v1alpha1.ClusterAddonSpec{
				Resources: []v1alpha1.RemoteRef{
					{
						OCI: v1alpha1.OCIRemoteRef{
							Name: cluster.GetEnv(cluster.EnvAddonDashboardOKDImage, imageDashboardOKD),
						},
					},
				},
				DisplayName: "OpenShift Console",
				Maintainer:  "RedHat",
				Logo:        "https://upload.wikimedia.org/wikipedia/commons/3/3a/OpenShift-LogoType.svg",
				Description: "The console is a more friendly kubectl in the form of a single page webapp (OpenShift-specific features such as Projects or Routes will not work, mutually-exclusive with the Kubernetes Dashboard).",
				Source:      v1alpha1.SourceCommunity,
				SourceURL:   "https://github.com/openshift/console",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "podinfo",
				Namespace: tr.GetName(),
				Labels:    labels,
			},
			Spec: v1alpha1.ClusterAddonSpec{
				Resources: []v1alpha1.RemoteRef{
					{
						OCI: v1alpha1.OCIRemoteRef{
							Name: cluster.GetEnv(cluster.EnvAddonPodInfoImage, imagePodInfo),
						},
					},
				},
				DisplayName: "PodInfo",
				Maintainer:  "KubeGlass",
				Logo:        "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif",
				Description: "Go microservice template for Kubernetes",
				Source:      v1alpha1.SourceOfficial,
				SourceURL:   "https://github.com/stefanprodan/podinfo",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "prometheus-adapter",
				Namespace: tr.GetName(),
				Labels:    labels,
			},
			Spec: v1alpha1.ClusterAddonSpec{
				Resources: []v1alpha1.RemoteRef{
					{
						OCI: v1alpha1.OCIRemoteRef{
							Name: cluster.GetEnv(cluster.EnvAddonMetricsAdapterImage, imagePromAdapter),
						},
					},
				},
				DisplayName: "Prometheus Adapter",
				Maintainer:  "KubeGlass",
				Logo:        "https://cncf-branding.netlify.app/img/projects/prometheus/icon/color/prometheus-icon-color.svg",
				Description: "Enables the Kubernetes Metrics API",
				Source:      v1alpha1.SourceOfficial,
				SourceURL:   "https://github.com/kubernetes-sigs/prometheus-adapter",
			},
		},
	}
}
