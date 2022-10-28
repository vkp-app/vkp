package cluster

import (
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1betav1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Cluster(cluster *paasv1alpha1.Cluster) *capiv1betav1.Cluster {
	clusterRef := &corev1.ObjectReference{
		Kind:       "VCluster",
		Namespace:  cluster.GetNamespace(),
		Name:       cluster.GetName(),
		APIVersion: "infrastructure.cluster.x-k8s.io/v1alpha1",
	}
	return &capiv1betav1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetName(),
			Namespace: cluster.GetNamespace(),
			Labels:    Labels(cluster),
		},
		Spec: capiv1betav1.ClusterSpec{
			ControlPlaneRef:   clusterRef,
			InfrastructureRef: clusterRef,
		},
	}
}
