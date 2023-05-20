package v1alpha1

import "k8s.io/apimachinery/pkg/util/validation/field"

const (
	KindAppliedClusterVersion = "AppliedClusterVersion"
	KindCluster               = "Cluster"
	KindClusterAddon          = "ClusterAddon"
	KindClusterAddonBinding   = "ClusterAddonBinding"
	KindClusterVersion        = "ClusterVersion"
	KindTenant                = "Tenant"
)

type nameMatchFunc = func(name string) *field.Error
