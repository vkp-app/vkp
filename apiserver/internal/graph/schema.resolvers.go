package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
)

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, input model.NewTenant) (*model.Tenant, error) {
	panic(fmt.Errorf("not implemented: CreateTenant - createTenant"))
}

// Tenants is the resolver for the tenants field.
func (r *queryResolver) Tenants(ctx context.Context) ([]*model.Tenant, error) {
	panic(fmt.Errorf("not implemented: Tenants - tenants"))
}

// Clusters is the resolver for the clusters field.
func (r *queryResolver) Clusters(ctx context.Context) ([]*model.Cluster, error) {
	panic(fmt.Errorf("not implemented: Clusters - clusters"))
}

// ClustersInTenant is the resolver for the clustersInTenant field.
func (r *queryResolver) ClustersInTenant(ctx context.Context, id string) ([]*model.Cluster, error) {
	panic(fmt.Errorf("not implemented: ClustersInTenant - clustersInTenant"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
