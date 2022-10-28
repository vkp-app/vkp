package cluster

import (
	"fmt"
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api/api/v1beta1"
)

func VCluster(cluster *paasv1alpha1.Cluster) *vclusterv1alpha1.VCluster {
	return &vclusterv1alpha1.VCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetName(),
			Namespace: cluster.GetNamespace(),
			Labels:    Labels(cluster),
		},
		Spec: vclusterv1alpha1.VClusterSpec{
			ControlPlaneEndpoint: v1beta1.APIEndpoint{
				// todo determine host in a better way
				Host: fmt.Sprintf("%s-%s.todo", cluster.GetName(), cluster.GetNamespace()),
				Port: 443,
			},
			HelmRelease: &vclusterv1alpha1.VirtualClusterHelmRelease{
				Chart: vclusterv1alpha1.VirtualClusterHelmChart{
					Name:    "vcluster",
					Repo:    "https://charts.loft.sh",
					Version: "0.12.2",
				},
				Values: "",
			},
			KubernetesVersion: pointer.String("1.24"),
		},
	}
}
