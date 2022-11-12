package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	cluster2 "gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cmdv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// Tenant is the resolver for the tenant field.
func (r *clusterResolver) Tenant(ctx context.Context, obj *paasv1alpha1.Cluster) (string, error) {
	return obj.ObjectMeta.Labels[labelTenant], nil
}

// Track is the resolver for the track field.
func (r *clusterResolver) Track(ctx context.Context, obj *paasv1alpha1.Cluster) (model.Track, error) {
	return model.FromDAO(obj.Spec.Track), nil
}

// DisplayName is the resolver for the displayName field.
func (r *clusterAddonResolver) DisplayName(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (string, error) {
	return obj.Spec.DisplayName, nil
}

// Description is the resolver for the description field.
func (r *clusterAddonResolver) Description(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (string, error) {
	return obj.Spec.Description, nil
}

// Maintainer is the resolver for the maintainer field.
func (r *clusterAddonResolver) Maintainer(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (string, error) {
	return obj.Spec.Maintainer, nil
}

// Logo is the resolver for the logo field.
func (r *clusterAddonResolver) Logo(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (string, error) {
	return obj.Spec.Logo, nil
}

// Source is the resolver for the source field.
func (r *clusterAddonResolver) Source(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (paasv1alpha1.AddonSource, error) {
	return obj.Spec.Source, nil
}

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, tenant string) (*paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("creating tenant")
	user, _ := userctx.CtxUser(ctx)
	// create the containing namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: tenant,
		},
	}
	if err := r.Create(ctx, ns); err != nil {
		log.Error(err, "failed to create tenant namespace")
		return nil, err
	}
	// create the tenant
	tr := &paasv1alpha1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tenant,
			Namespace: tenant,
		},
		Spec: paasv1alpha1.TenantSpec{
			Owner:             user.Username,
			NamespaceStrategy: paasv1alpha1.StrategySingle,
		},
		Status: paasv1alpha1.TenantStatus{
			// all tenants must be approved by either
			// a human administrator or some sort of automated approval
			// (e.g. payment-method verification)
			Phase: paasv1alpha1.PhasePendingApproval,
		},
	}
	if err := r.Create(ctx, tr); err != nil {
		log.Error(err, "failed to create tenant")
		return nil, err
	}
	// ensure that the correct phase is applied
	if err := r.Status().Update(ctx, tr); err != nil {
		log.Error(err, "failed to apply tenant status")
		return nil, err
	}
	return tr, nil
}

// CreateCluster is the resolver for the createCluster field.
func (r *mutationResolver) CreateCluster(ctx context.Context, tenant string, input model.NewCluster) (*paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("creating cluster")
	// validate the tenant
	tenantResource := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: tenant}, tenantResource); err != nil {
		log.Error(err, "failed to retrieve tenant information")
		return nil, err
	}
	// reject clusters for tenants that have yet
	// to be approved
	if tenantResource.Status.Phase != paasv1alpha1.PhaseReady {
		log.Info("rejecting cluster creation request for tenant that is not 'Ready'")
		return nil, ErrTenantNotReady
	}

	// create the cluster
	cluster := &paasv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: tenant,
			Labels: map[string]string{
				labelTenant: tenant,
			},
		},
		Spec: paasv1alpha1.ClusterSpec{
			Track: input.Track.ToDAO(),
		},
	}
	if err := r.Create(ctx, cluster); err != nil {
		log.Error(err, "failed to create cluster")
		return nil, err
	}
	return cluster, nil
}

// ApproveTenant is the resolver for the approveTenant field.
func (r *mutationResolver) ApproveTenant(ctx context.Context, tenant string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("approving tenant")
	// validate the tenant
	tenantResource := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: tenant}, tenantResource); err != nil {
		log.Error(err, "failed to retrieve tenant information")
		return false, err
	}
	// if the tenant doesn't require approval
	// exit cleanly
	if tenantResource.Status.Phase != paasv1alpha1.PhasePendingApproval {
		log.Info("tenant is not awaiting approval")
		return false, nil
	}
	// update the .status.phase field to indicate
	// that the tenant is ready for use
	tenantResource.Status.Phase = paasv1alpha1.PhaseReady
	if err := r.Status().Update(ctx, tenantResource); err != nil {
		log.Error(err, "failed to update tenant phase")
		return false, err
	}
	return true, nil
}

// Tenants is the resolver for the tenants field.
func (r *queryResolver) Tenants(ctx context.Context) ([]paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("listing tenants")
	tenants := &paasv1alpha1.TenantList{}
	if err := r.List(ctx, tenants); err != nil {
		log.Error(err, "failed to list tenants")
		return nil, err
	}
	return tenants.Items, nil
}

// Tenant is the resolver for the tenant field.
func (r *queryResolver) Tenant(ctx context.Context, tenant string) (*paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("fetching tenant")
	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: tenant}, tr); err != nil {
		log.Error(err, "failed to get tenant")
		return nil, err
	}
	return tr, nil
}

// ClustersInTenant is the resolver for the clustersInTenant field.
func (r *queryResolver) ClustersInTenant(ctx context.Context, tenant string) ([]paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("fetching clusters in tenant")
	clusters := &paasv1alpha1.ClusterList{}
	if err := r.List(ctx, clusters, client.MatchingLabels{labelTenant: tenant}); err != nil {
		log.Error(err, "failed to list clusters in tenant")
		return nil, err
	}
	return clusters.Items, nil
}

// Cluster is the resolver for the cluster field.
func (r *queryResolver) Cluster(ctx context.Context, tenant string, cluster string) (*paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("fetching cluster")
	cr := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: cluster}, cr); err != nil {
		log.Error(err, "failed to retrieve cluster")
		return nil, err
	}
	return cr, nil
}

// ClusterAddons is the resolver for the clusterAddons field.
func (r *queryResolver) ClusterAddons(ctx context.Context, tenant string) ([]paasv1alpha1.ClusterAddon, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("listing addons")

	addons := &paasv1alpha1.ClusterAddonList{}
	if err := r.List(ctx, addons, client.InNamespace(tenant)); err != nil {
		log.Error(err, "failed to list addons")
		return nil, err
	}
	return addons.Items, nil
}

// ClusterInstalledAddons is the resolver for the clusterInstalledAddons field.
func (r *queryResolver) ClusterInstalledAddons(ctx context.Context, tenant string, cluster string) ([]string, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("listing installed addons")

	addons := &paasv1alpha1.ClusterAddonBindingList{}
	if err := r.List(ctx, addons, client.MatchingLabels{paasv1alpha1.LabelClusterRef: cluster}, client.InNamespace(tenant)); err != nil {
		log.Error(err, "failed to list addon bindings")
	}
	// collect the list of names
	names := make([]string, len(addons.Items))
	for i := range addons.Items {
		names[i] = addons.Items[i].GetName()
	}
	return names, nil
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user, _ := userctx.CtxUser(ctx)
	return user, nil
}

// ClusterMetricMemory is the resolver for the clusterMetricMemory field.
func (r *queryResolver) ClusterMetricMemory(ctx context.Context, tenant string, cluster string) ([]model.MetricValue, error) {
	return r.GetMetric(ctx, fmt.Sprintf(`sum by (namespace) (container_memory_usage_bytes{namespace="%s", pod=~".*-%s|%s-.+"})`, tenant, cluster, cluster))
}

// ClusterMetricCPU is the resolver for the clusterMetricCPU field.
func (r *queryResolver) ClusterMetricCPU(ctx context.Context, tenant string, cluster string) ([]model.MetricValue, error) {
	return r.GetMetric(ctx, fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s", pod=~".*-%s|%s-.+"}[1m])) by (namespace)`, tenant, cluster, cluster))
}

// ClusterMetricPods is the resolver for the clusterMetricPods field.
func (r *queryResolver) ClusterMetricPods(ctx context.Context, tenant string, cluster string) ([]model.MetricValue, error) {
	return r.GetMetric(ctx, fmt.Sprintf(`sum by (namespace) (kube_pod_status_ready{namespace="%s", pod=~".*-%s|%s-.+", condition="true"})`, tenant, cluster, cluster))
}

// ClusterMetricNetReceive is the resolver for the clusterMetricNetReceive field.
func (r *queryResolver) ClusterMetricNetReceive(ctx context.Context, tenant string, cluster string) ([]model.MetricValue, error) {
	return r.GetMetric(ctx, fmt.Sprintf(`sum by (namespace) (irate(node_network_receive_bytes_total{namespace="%s", pod=~".*-%s|%s-.+"}[2m]))`, tenant, cluster, cluster))
}

// ClusterMetricNetTransmit is the resolver for the clusterMetricNetTransmit field.
func (r *queryResolver) ClusterMetricNetTransmit(ctx context.Context, tenant string, cluster string) ([]model.MetricValue, error) {
	return r.GetMetric(ctx, fmt.Sprintf(`sum by (namespace) (irate(node_network_transmit_bytes_total{namespace="%s", pod=~".*-%s|%s-.+"}[2m]))`, tenant, cluster, cluster))
}

// RenderKubeconfig is the resolver for the renderKubeconfig field.
func (r *queryResolver) RenderKubeconfig(ctx context.Context, tenant string, cluster string) (string, error) {
	log := logr.FromContextOrDiscard(ctx)
	user, _ := userctx.CtxUser(ctx)

	// fetch the cluster resource
	cr, err := r.Cluster(ctx, tenant, cluster)
	if err != nil {
		return "", err
	}

	// fetch the dex secret for this cluster
	dexSec := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cluster2.DexSecretName(cr.GetName())}, dexSec); err != nil {
		log.Error(err, "failed to retrieve dex secret")
		return "", err
	}

	// fetch the TLS secret for this cluster
	tlsSec := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cluster2.IngressSecretName(cr.GetName())}, tlsSec); err != nil {
		log.Error(err, "failed to retrieve TLS secret")
		return "", err
	}

	username := strings.Split(user.Username, "@")[0]
	clusterName := fmt.Sprintf("%s-%s", tenant, cluster)
	contextName := fmt.Sprintf("%s@%s", username, clusterName)

	// create the kubeconfig struct and
	// populate the information that the user
	// will need
	cfg := cmdv1.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		CurrentContext: contextName,
		Clusters: []cmdv1.NamedCluster{
			{
				Name: clusterName,
				Cluster: cmdv1.Cluster{
					Server:                   fmt.Sprintf("https://%s:443", cr.Status.KubeURL),
					CertificateAuthorityData: tlsSec.Data["ca.crt"],
				},
			},
		},
		Contexts: []cmdv1.NamedContext{
			{
				Name: contextName,
				Context: cmdv1.Context{
					Cluster:   clusterName,
					AuthInfo:  username,
					Namespace: "default",
				},
			},
		},
		AuthInfos: []cmdv1.NamedAuthInfo{
			{
				Name: username,
				AuthInfo: cmdv1.AuthInfo{
					Exec: &cmdv1.ExecConfig{
						APIVersion: "client.authentication.k8s.io/v1beta1",
						Command:    "kubectl",
						Args: []string{
							"oidc-login",
							"get-token",
							fmt.Sprintf("--oidc-issuer-url=%s", os.Getenv(cluster2.EnvIDPURL)),
							fmt.Sprintf("--oidc-client-id=%s", string(dexSec.Data[cluster2.DexKeyID])),
							fmt.Sprintf("--oidc-client-secret=%s", string(dexSec.Data[cluster2.DexKeySecret])),
							"--oidc-extra-scope=profile",
							"--oidc-extra-scope=email",
							"--oidc-extra-scope=groups",
							fmt.Sprintf("--certificate-authority-data=%s", base64.StdEncoding.EncodeToString(tlsSec.Data["ca.crt"])),
						},
						InstallHint: "kubectl krew install oidc-login",
					},
				},
			},
		},
	}
	// convert the struct to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Error(err, "failed to convert cmdv1.Config to YAML")
		return "", err
	}
	return string(data), nil
}

// Owner is the resolver for the owner field.
func (r *tenantResolver) Owner(ctx context.Context, obj *paasv1alpha1.Tenant) (string, error) {
	return obj.Spec.Owner, nil
}

// ObservedClusters is the resolver for the observedClusters field.
func (r *tenantResolver) ObservedClusters(ctx context.Context, obj *paasv1alpha1.Tenant) ([]paasv1alpha1.NamespacedName, error) {
	return obj.Status.ObservedClusters, nil
}

// Cluster returns generated.ClusterResolver implementation.
func (r *Resolver) Cluster() generated.ClusterResolver { return &clusterResolver{r} }

// ClusterAddon returns generated.ClusterAddonResolver implementation.
func (r *Resolver) ClusterAddon() generated.ClusterAddonResolver { return &clusterAddonResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Tenant returns generated.TenantResolver implementation.
func (r *Resolver) Tenant() generated.TenantResolver { return &tenantResolver{r} }

type clusterResolver struct{ *Resolver }
type clusterAddonResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type tenantResolver struct{ *Resolver }
