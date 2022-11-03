package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Tenant is the resolver for the tenant field.
func (r *clusterResolver) Tenant(ctx context.Context, obj *paasv1alpha1.Cluster) (string, error) {
	return obj.ObjectMeta.Labels[labelTenant], nil
}

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, name string) (*paasv1alpha1.Tenant, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", name)
	log.Info("creating tenant")
	user, ok := userctx.CtxUser(ctx)
	if !ok {
		return nil, ErrUnauthorised
	}
	// create the containing namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := r.Create(ctx, ns); err != nil {
		log.Error(err, "failed to create tenant namespace")
		return nil, err
	}
	// create the tenant
	tenant := &paasv1alpha1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: name,
		},
		Spec: paasv1alpha1.TenantSpec{
			Owner:             user.Username,
			NamespaceStrategy: paasv1alpha1.StrategySingle,
		},
	}
	if err := r.Create(ctx, tenant); err != nil {
		log.Error(err, "failed to create tenant")
		return nil, err
	}
	return tenant, nil
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

// ClustersInTenant is the resolver for the clustersInTenant field.
func (r *queryResolver) ClustersInTenant(ctx context.Context, tenant string) ([]paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant)
	log.Info("fetching clusters in tenant")
	clusters := &paasv1alpha1.ClusterList{}
	selector := labels.SelectorFromSet(labels.Set{labelTenant: tenant})
	if err := r.List(ctx, clusters, &client.ListOptions{LabelSelector: selector}); err != nil {
		log.Error(err, "failed to list clusters in tenant")
		return nil, err
	}
	return clusters.Items, nil
}

// Cluster is the resolver for the cluster field.
func (r *queryResolver) Cluster(ctx context.Context, tenant string, name string) (*paasv1alpha1.Cluster, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", name)
	log.Info("fetching cluster")
	cluster := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: name}, cluster); err != nil {
		log.Error(err, "failed to retrieve cluster")
		return nil, err
	}
	panic(fmt.Errorf("not implemented: Cluster - cluster"))
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user, ok := userctx.CtxUser(ctx)
	if !ok {
		return nil, ErrUnauthorised
	}
	return user, nil
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Tenant returns generated.TenantResolver implementation.
func (r *Resolver) Tenant() generated.TenantResolver { return &tenantResolver{r} }

type clusterResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type tenantResolver struct{ *Resolver }
