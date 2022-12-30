package cluster

import (
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1betav1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	vclusterKind       = "VCluster"
	vclusterApiVersion = "infrastructure.cluster.x-k8s.io/v1alpha1"
)

func Cluster(cluster *paasv1alpha1.Cluster) *capiv1betav1.Cluster {
	clusterRef := &corev1.ObjectReference{
		Kind:       vclusterKind,
		Namespace:  cluster.GetNamespace(),
		Name:       cluster.GetName(),
		APIVersion: vclusterApiVersion,
	}
	return &capiv1betav1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetName(),
			Namespace: cluster.GetNamespace(),
			Labels:    Labels(cluster),
		},
		Spec: capiv1betav1.ClusterSpec{
			ControlPlaneEndpoint: capiv1betav1.APIEndpoint{
				Host: getHostname(cluster),
				Port: 443,
			},
			ControlPlaneRef:   clusterRef,
			InfrastructureRef: clusterRef,
		},
	}
}

func AppliedClusterVersion(cr *paasv1alpha1.Cluster) *paasv1alpha1.AppliedClusterVersion {
	return &paasv1alpha1.AppliedClusterVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    Labels(cr),
		},
		Spec: paasv1alpha1.AppliedClusterVersionSpec{
			ClusterRef: corev1.LocalObjectReference{
				Name: cr.GetName(),
			},
		},
	}
}
