package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/djcass44/go-utils/utilities/sliceutils"
	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

func (r *Resolver) HasTenantAccess(ctx context.Context, req any, next graphql.Resolver, write bool) (res any, err error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("preparing to check tenant access", "request", req)

	// chose not to do an "args, ok := ..." check here
	// since it wouldn't show us the interface conversion error
	args := req.(map[string]any)
	tenant := args["tenant"].(string)

	ok, err := r.canAccessTenant(ctx, tenant, write)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}

	return next(ctx)
}

func (r *Resolver) HasClusterAccess(ctx context.Context, req any, next graphql.Resolver, write bool) (res any, err error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("preparing to check cluster access", "request", req)

	// chose not to do an "args, ok := ..." check here
	// since it wouldn't show us the interface conversion error
	args := req.(map[string]any)
	cluster := args["cluster"].(string)
	tenant := args["tenant"].(string)

	ok, err := r.canAccessCluster(ctx, tenant, cluster, write)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}

	return next(ctx)
}

func (r *Resolver) canAccessTenant(ctx context.Context, tenant string, requiresWrite bool) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "requiresWrite", requiresWrite)

	// fetch the tenant resource
	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tr); err != nil {
		log.Error(err, "failed to retrieve tenant resource")
		return false, err
	}

	return r.canAccessTenantResource(ctx, tr, requiresWrite)
}

func (r *Resolver) canAccessTenantResource(ctx context.Context, tr *paasv1alpha1.Tenant, requiresWrite bool) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tr.GetName(), "requiresWrite", requiresWrite)
	log.Info("checking tenant access")
	user, _ := userctx.CtxUser(ctx)

	// if the user owns the tenant, then
	// that's it really.
	if tr.Spec.Owner == user.Username {
		log.V(1).Info("user can access tenant as they own it")
		return true, nil
	}

	for _, ar := range tr.Spec.Accessors {
		if ar.User == user.Username || sliceutils.Includes(user.Groups, ar.Group) {
			// if this request requires write access
			// ignore all accessors that provide
			// read-only access
			if requiresWrite || ar.ReadOnly {
				continue
			}
			log.V(1).Info("user can access the tenant as they are referenced in an accessor")
			return true, nil
		}
	}

	log.V(1).Info("user cannot access tenant as we failed to locate a matching accessor")
	return false, nil
}

func (r *Resolver) canAccessCluster(ctx context.Context, tenant, cluster string, requiresWrite bool) (bool, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster, "requiresWrite", requiresWrite)
	log.Info("checking cluster access")

	user, _ := userctx.CtxUser(ctx)

	// fetch the tenant since we can
	// short-circuit the request if the user
	// owns it
	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, types.NamespacedName{Name: tenant}, tr); err != nil {
		log.Error(err, "failed to retrieve tenant resource")
		return false, err
	}
	// if the user owns the tenant, they
	// own the cluster.
	if tr.Spec.Owner == user.Username {
		log.V(1).Info("user can access cluster as they own the tenant")
		return true, nil
	}

	// fetch the cluster resource
	cr := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tenant, Name: cluster}, cr); err != nil {
		log.Error(err, "failed to retrieve cluster resource")
		return false, err
	}

	for _, ar := range cr.Spec.Accessors {
		if ar.User == user.Username || sliceutils.Includes(user.Groups, ar.Group) {
			// if this request requires write access
			// ignore all accessors that provide
			// read-only access
			if requiresWrite || ar.ReadOnly {
				continue
			}
			log.V(1).Info("user can access the cluster as they are referenced in an accessor")
			return true, nil
		}
	}

	log.V(1).Info("user cannot access cluster as we failed to locate a matching accessor")
	return false, nil
}

func (r *Resolver) HasRole(ctx context.Context, _ any, next graphql.Resolver, role model.Role) (any, error) {
	log := logr.FromContextOrDiscard(ctx)
	user, ok := userctx.CtxUser(ctx)
	if !ok {
		log.V(1).Info("rejecting unauthorised request")
		return nil, ErrUnauthorised
	}
	// skip the admin check if this endpoint
	// only need user privilege
	if role == model.RoleUser {
		return next(ctx)
	}
	if err := r.userHasAdmin(ctx, user); err != nil {
		return nil, err
	}
	return next(ctx)
}

func (r *Resolver) userHasAdmin(ctx context.Context, user *model.User) error {
	log := logr.FromContextOrDiscard(ctx)
	user, ok := userctx.CtxUser(ctx)
	if !ok {
		return ErrUnauthorised
	}
	// run a SAR to verify that the user can access management-cluster resources
	log.Info("verifying administrative privileges of requesting user")
	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				// if the user can get secrets in the glass namespace
				// then they can be considered an administrator
				// of the system
				Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
				Verb:      "get",
				Version:   "v1",
				Resource:  "secrets",
			},
			User:   user.Username,
			Groups: user.Groups,
		},
	}
	if err := r.Create(ctx, sar); err != nil {
		log.Error(err, "failed to create SubjectAccessReview")
		return err
	}
	if !sar.Status.Allowed {
		log.Info("rejecting admin request due to failing SAR checks")
		return ErrForbidden
	}
	return nil
}
