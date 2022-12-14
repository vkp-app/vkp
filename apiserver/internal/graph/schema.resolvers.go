package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/robfig/cron"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
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
func (r *clusterResolver) Track(ctx context.Context, obj *paasv1alpha1.Cluster) (paasv1alpha1.ReleaseTrack, error) {
	return obj.Spec.Track, nil
}

// Accessors is the resolver for the accessors field.
func (r *clusterResolver) Accessors(ctx context.Context, obj *paasv1alpha1.Cluster) ([]paasv1alpha1.AccessRef, error) {
	return obj.Spec.Accessors, nil
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

// SourceURL is the resolver for the sourceURL field.
func (r *clusterAddonResolver) SourceURL(ctx context.Context, obj *paasv1alpha1.ClusterAddon) (string, error) {
	return obj.Spec.SourceURL, nil
}

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, tenant string) (*paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("creating tenant")
	user, _ := userctx.CtxUser(ctx)
	// create the tenant
	tr := &paasv1alpha1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name: tenant,
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
	// ensure that the correct phase is applied.
	// this may fail in certain situations (e.g. OpenShift) so
	// don't error on failure
	if err := r.Status().Update(ctx, tr); err != nil {
		log.Error(err, "failed to apply tenant status")
		return tr, nil
	}
	return tr, nil
}

// CreateCluster is the resolver for the createCluster field.
func (r *mutationResolver) CreateCluster(ctx context.Context, tenant string, input model.NewCluster) (*paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("creating cluster")
	user, _ := userctx.CtxUser(ctx)
	// validate the tenant
	tenantResource := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tenantResource); err != nil {
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
			Track: input.Track,
			HA: paasv1alpha1.HighAvailability{
				Enabled: input.Ha,
			},
			Accessors: []paasv1alpha1.AccessRef{
				{
					// make sure that the creating user has administrative privileges
					// over the cluster
					User: user.Username,
				},
			},
		},
	}
	if err := r.Create(ctx, cluster); err != nil {
		log.Error(err, "failed to create cluster")
		return nil, err
	}
	return cluster, nil
}

// DeleteCluster is the resolver for the deleteCluster field.
func (r *mutationResolver) DeleteCluster(ctx context.Context, tenant string, cluster string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("deleting cluster")

	// delete all the addons that were
	// installed for this cluster
	if err := r.DeleteAllOf(ctx, &paasv1alpha1.ClusterAddonBinding{}, client.InNamespace(tenant), client.MatchingLabels{paasv1alpha1.LabelClusterRef: cluster}); err != nil {
		log.Error(err, "failed to delete cluster addons")
		return false, err
	}

	// create a dummy cluster object, so we
	// can delete the cluster
	cr := &paasv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster,
			Namespace: tenant,
			Labels: map[string]string{
				labelTenant: tenant,
			},
		},
	}
	if err := r.Delete(ctx, cr); err != nil {
		log.Error(err, "failed to delete cluster")
		return false, err
	}

	return true, nil
}

// SetClusterAccessors is the resolver for the setClusterAccessors field.
func (r *mutationResolver) SetClusterAccessors(ctx context.Context, tenant string, cluster string, accessors []model.AccessRefInput) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("updating cluster accessors")

	// fetch the cluster
	cr := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: cluster}, cr); err != nil {
		log.Error(err, "failed to retrieve cluster resource")
		return false, err
	}
	newAccess := make([]paasv1alpha1.AccessRef, len(accessors))
	for i := range accessors {
		newAccess[i] = paasv1alpha1.AccessRef{
			ReadOnly: accessors[i].ReadOnly,
			User:     accessors[i].User,
			Group:    accessors[i].Group,
		}
	}
	log.V(1).Info("applying new accessors", "old", cr.Spec.Accessors, "new", newAccess)
	cr.Spec.Accessors = newAccess
	// save changes
	if err := r.Update(ctx, cr); err != nil {
		log.Error(err, "failed to update cluster")
		return false, err
	}
	return true, nil
}

// SetClusterMaintenanceWindow is the resolver for the setClusterMaintenanceWindow field.
func (r *mutationResolver) SetClusterMaintenanceWindow(ctx context.Context, tenant string, cluster string, window string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("updating cluster maintenance window")
	// fetch the AppliedClusterVersion
	acv := &paasv1alpha1.AppliedClusterVersion{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: cluster}, acv); err != nil {
		log.Error(err, "failed to retrieve AppliedClusterVersion")
		return false, fmt.Errorf("failed to retrieve cluster version information: %w", err)
	}

	// update the resource
	acv.Spec.MaintenanceWindow = window
	if err := r.Update(ctx, acv); err != nil {
		log.Error(err, "failed to update AppliedClusterVersion")
		return false, fmt.Errorf("failed to update maintenance window: %w", err)
	}
	return true, nil
}

// InstallAddon is the resolver for the installAddon field.
func (r *mutationResolver) InstallAddon(ctx context.Context, tenant string, cluster string, addon string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster, "addon", addon)
	log.Info("installing addon")

	// todo validate that the cluster is ready

	// create the binding
	br := &paasv1alpha1.ClusterAddonBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", cluster, addon),
			Namespace: tenant,
		},
		Spec: paasv1alpha1.ClusterAddonBindingSpec{
			ClusterRef: corev1.LocalObjectReference{
				Name: cluster,
			},
			ClusterAddonRef: corev1.LocalObjectReference{
				Name: addon,
			},
		},
	}
	if err := r.Create(ctx, br); err != nil {
		log.Error(err, "failed to create ClusterAddonBinding")
		return false, err
	}
	br.Status.Phase = paasv1alpha1.AddonPhaseInstalled
	if err := r.Status().Update(ctx, br); err != nil {
		log.Error(err, "failed to update ClusterAddonBinding status")
		return false, err
	}
	return true, nil
}

// UninstallAddon is the resolver for the uninstallAddon field.
func (r *mutationResolver) UninstallAddon(ctx context.Context, tenant string, cluster string, addon string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster, "addon", addon)
	log.Info("uninstalling addon")

	// create a shell
	br := &paasv1alpha1.ClusterAddonBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", cluster, addon),
			Namespace: tenant,
		},
	}
	if err := r.Delete(ctx, br); err != nil {
		log.Error(err, "failed to delete ClusterAddonBinding")
		return false, err
	}
	log.Info("completed addon uninstallation")
	return true, nil
}

// ApproveTenant is the resolver for the approveTenant field.
func (r *mutationResolver) ApproveTenant(ctx context.Context, tenant string) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("approving tenant")
	// validate the tenant
	tenantResource := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tenantResource); err != nil {
		log.Error(err, "failed to retrieve tenant information")
		return false, err
	}
	// if the tenant doesn't require approval
	// exit cleanly
	if tenantResource.Status.Phase != paasv1alpha1.PhasePendingApproval && paasv1alpha1.PhasePendingApproval != "" {
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

// SetTenantAccessors is the resolver for the setTenantAccessors field.
func (r *mutationResolver) SetTenantAccessors(ctx context.Context, tenant string, accessors []model.AccessRefInput) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("updating tenant accessors")

	// fetch the tenant
	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tr); err != nil {
		log.Error(err, "failed to retrieve cluster resource")
		return false, err
	}
	newAccess := make([]paasv1alpha1.AccessRef, len(accessors))
	for i := range accessors {
		groupName := accessors[i].Group
		// add the 'oidc:' prefix if
		// the user doesn't specify it
		// themselves
		if groupName != "" && !strings.HasPrefix(groupName, "oidc:") {
			groupName = "oidc:" + groupName
		}
		newAccess[i] = paasv1alpha1.AccessRef{
			ReadOnly: accessors[i].ReadOnly,
			User:     accessors[i].User,
			Group:    groupName,
		}
	}
	log.V(1).Info("applying new accessors", "old", tr.Spec.Accessors, "new", newAccess)
	tr.Spec.Accessors = newAccess
	// save changes
	if err := r.Update(ctx, tr); err != nil {
		log.Error(err, "failed to update tenant")
		return false, err
	}
	return true, nil
}

// Tenants is the resolver for the tenants field.
func (r *queryResolver) Tenants(ctx context.Context) ([]paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("listing tenants")
	user, _ := userctx.CtxUser(ctx)

	tenants := &paasv1alpha1.TenantList{}
	if err := r.List(ctx, tenants); err != nil {
		log.Error(err, "failed to list tenants")
		return nil, err
	}
	// if the user is an admin, return the full list
	if err := r.userHasAdmin(ctx, user); err == nil {
		return tenants.Items, nil
	}

	// otherwise, filter the tenants
	var items []paasv1alpha1.Tenant
	for _, tr := range tenants.Items {
		// todo figure out a better way of filtering tenants
		// since this will probably have a pretty
		// nasty performance penalty
		ok, err := r.canAccessTenantResource(ctx, &tr, false)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		items = append(items, tr)
	}
	return items, nil
}

// Tenant is the resolver for the tenant field.
func (r *queryResolver) Tenant(ctx context.Context, tenant string) (*paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("fetching tenant")
	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tr); err != nil {
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
func (r *queryResolver) ClusterInstalledAddons(ctx context.Context, tenant string, cluster string) ([]model.AddonBindingStatus, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("listing installed addons")

	addons := &paasv1alpha1.ClusterAddonBindingList{}
	if err := r.List(ctx, addons, client.MatchingLabels{paasv1alpha1.LabelClusterRef: cluster}, client.InNamespace(tenant)); err != nil {
		log.Error(err, "failed to list addon bindings")
	}
	// collect the list of names
	//goland:noinspection GoPreferNilSlice
	names := []model.AddonBindingStatus{}
	for i := range addons.Items {
		names = append(names, model.AddonBindingStatus{
			Name:  addons.Items[i].GetName(),
			Phase: addons.Items[i].Status.Phase,
		})
	}
	return names, nil
}

// ClusterMaintenanceWindow is the resolver for the clusterMaintenanceWindow field.
func (r *queryResolver) ClusterMaintenanceWindow(ctx context.Context, tenant string, cluster string) (*model.MaintenanceWindow, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)
	log.Info("updating cluster maintenance window")
	// fetch the AppliedClusterVersion
	acv := &paasv1alpha1.AppliedClusterVersion{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: cluster}, acv); err != nil {
		log.Error(err, "failed to retrieve AppliedClusterVersion")
		return nil, fmt.Errorf("failed to retrieve cluster version information: %w", err)
	}
	schedule, err := cron.ParseStandard(acv.Spec.MaintenanceWindow)
	if err != nil {
		log.Error(err, "failed to parse maintenance window schedule - this should never happen")
		return nil, err
	}
	return &model.MaintenanceWindow{
		Schedule: acv.Spec.MaintenanceWindow,
		Next:     schedule.Next(time.Now()).Unix(),
	}, nil
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user, _ := userctx.CtxUser(ctx)
	return user, nil
}

// ClusterMetrics is the resolver for the clusterMetrics field.
func (r *queryResolver) ClusterMetrics(ctx context.Context, tenant string, cluster string) ([]model.Metric, error) {
	return r.GetMetrics(ctx, tenant, cluster)
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

	// figure out which CA we need to use
	var caData []byte
	// use the kube API CA if
	// we don't have one for Dex.
	if r.dexCA == "" {
		caData = tlsSec.Data["ca.crt"]
	} else {
		caData = []byte(r.dexCA)
	}

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
							fmt.Sprintf("--oidc-issuer-url=%s", r.dexURL),
							fmt.Sprintf("--oidc-client-id=%s", string(dexSec.Data[cluster2.DexKeyID])),
							fmt.Sprintf("--oidc-client-secret=%s", string(dexSec.Data[cluster2.DexKeySecret])),
							"--oidc-extra-scope=profile",
							"--oidc-extra-scope=email",
							"--oidc-extra-scope=groups",
							fmt.Sprintf("--certificate-authority-data=%s", base64.StdEncoding.EncodeToString(caData)),
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

// HasRole is the resolver for the hasRole field.
func (r *queryResolver) HasRole(ctx context.Context, role model.Role) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("role", role.String())
	log.Info("checking if user has role")
	user, _ := userctx.CtxUser(ctx)
	switch role {
	case model.RoleUser:
		return true, nil
	case model.RoleAdmin:
		err := r.userHasAdmin(ctx, user)
		return err == nil, err
	default:
		return false, fmt.Errorf("unknown role: %s", role.String())
	}
}

// HasTenantAccess is the resolver for the hasTenantAccess field.
func (r *queryResolver) HasTenantAccess(ctx context.Context, tenant string, write bool) (bool, error) {
	return r.canAccessTenant(ctx, tenant, write)
}

// HasClusterAccess is the resolver for the hasClusterAccess field.
func (r *queryResolver) HasClusterAccess(ctx context.Context, tenant string, cluster string, write bool) (bool, error) {
	return r.canAccessCluster(ctx, tenant, cluster, write)
}

// Owner is the resolver for the owner field.
func (r *tenantResolver) Owner(ctx context.Context, obj *paasv1alpha1.Tenant) (string, error) {
	return obj.Spec.Owner, nil
}

// ObservedClusters is the resolver for the observedClusters field.
func (r *tenantResolver) ObservedClusters(ctx context.Context, obj *paasv1alpha1.Tenant) ([]paasv1alpha1.NamespacedName, error) {
	return obj.Status.ObservedClusters, nil
}

// Accessors is the resolver for the accessors field.
func (r *tenantResolver) Accessors(ctx context.Context, obj *paasv1alpha1.Tenant) ([]paasv1alpha1.AccessRef, error) {
	return obj.Spec.Accessors, nil
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
