package cluster

import (
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Addons(tr *paasv1alpha1.Tenant) []paasv1alpha1.ClusterAddon {
	labels := TenantLabels(tr)
	return []paasv1alpha1.ClusterAddon{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dashboard-k8s",
				Namespace: tr.GetNamespace(),
				Labels:    labels,
			},
			Spec: paasv1alpha1.ClusterAddonSpec{
				Resources: []paasv1alpha1.RemoteRef{
					{
						OCI: "harbor.dcas.dev/kube-glass/addons/dashboard-k8s:1.0.0",
					},
				},
				DisplayName: "Kubernetes Dashboard",
				Maintainer:  "The Kubernetes Authors",
				Logo:        "https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo.svg",
				Description: "General-purpose web UI for Kubernetes clusters (Mutually-exclusive with the OpenShift Console).",
				Source:      paasv1alpha1.SourceCommunity,
				SourceURL:   "https://github.com/kubernetes/dashboard",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dashboard-okd",
				Namespace: tr.GetNamespace(),
				Labels:    labels,
			},
			Spec: paasv1alpha1.ClusterAddonSpec{
				Resources: []paasv1alpha1.RemoteRef{
					{
						OCI: "harbor.dcas.dev/kube-glass/addons/dashboard-okd:1.0.0",
					},
				},
				DisplayName: "OpenShift Console",
				Maintainer:  "RedHat",
				Logo:        "https://upload.wikimedia.org/wikipedia/commons/3/3a/OpenShift-LogoType.svg",
				Description: "The console is a more friendly kubectl in the form of a single page webapp (OpenShift-specific features such as Projects or Routes will not work, mutually-exclusive with the Kubernetes Dashboard).",
				Source:      paasv1alpha1.SourceCommunity,
				SourceURL:   "https://github.com/openshift/console",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "podinfo",
				Namespace: tr.GetNamespace(),
				Labels:    labels,
			},
			Spec: paasv1alpha1.ClusterAddonSpec{
				Resources: []paasv1alpha1.RemoteRef{
					{
						OCI: "harbor.dcas.dev/kube-glass/addons/podinfo:1.0.0",
					},
				},
				DisplayName: "PodInfo",
				Maintainer:  "KubeGlass",
				Logo:        "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif",
				Description: "Go microservice template for Kubernetes",
				Source:      paasv1alpha1.SourceOfficial,
				SourceURL:   "https://github.com/stefanprodan/podinfo",
			},
		},
	}
}
