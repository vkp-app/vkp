package nested

import (
	"fmt"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	authv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	rbacAPIGroup    = "rbac.authorization.k8s.io"
	kindClusterRole = "ClusterRole"

	roleClusterAdmin = "cluster-admin"
	roleView         = "view"
)

func OwnerBinding(cr *paasv1alpha1.Tenant) *authv1.ClusterRoleBinding {
	return &authv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-admin-owner",
		},
		RoleRef: authv1.RoleRef{
			Name:     roleClusterAdmin,
			Kind:     kindClusterRole,
			APIGroup: rbacAPIGroup,
		},
		Subjects: []authv1.Subject{
			{
				Kind:     authv1.UserKind,
				APIGroup: rbacAPIGroup,
				Name:     cr.Spec.Owner,
			},
		},
	}
}

func AccessWriteBinding(cr *paasv1alpha1.AccessRef) *authv1.ClusterRoleBinding {
	name := cr.Group
	kind := authv1.GroupKind
	if cr.User != "" {
		kind = authv1.UserKind
		name = cr.User
	}
	role := roleClusterAdmin
	if cr.ReadOnly {
		role = roleView
	}
	// we need to generate a name that we can safely store
	// in kubernetes that contains all the information
	// required to deterministically fetch this resource
	// at a later time.
	resourceName := fmt.Sprintf("accessor-%s-%s-%s", strings.ToLower(kind), name, role)
	return &authv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: idToName(resourceName),
		},
		RoleRef: authv1.RoleRef{
			Name:     role,
			Kind:     kindClusterRole,
			APIGroup: rbacAPIGroup,
		},
		Subjects: []authv1.Subject{
			{
				Kind:     kind,
				APIGroup: rbacAPIGroup,
				Name:     name,
			},
		},
	}
}
