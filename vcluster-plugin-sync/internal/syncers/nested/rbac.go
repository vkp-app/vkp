package nested

import (
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	authv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	rbacAPIGroup    = "rbac.authorization.k8s.io"
	kindClusterRole = "ClusterRole"
	kindUser        = "User"
)

func OwnerBinding(cr *paasv1alpha1.Tenant) *authv1.ClusterRoleBinding {
	return &authv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-admin-owner",
		},
		RoleRef: authv1.RoleRef{
			Name:     "cluster-admin",
			Kind:     kindClusterRole,
			APIGroup: rbacAPIGroup,
		},
		Subjects: []authv1.Subject{
			{
				Kind:     kindUser,
				APIGroup: rbacAPIGroup,
				Name:     cr.Spec.Owner,
			},
		},
	}
}
