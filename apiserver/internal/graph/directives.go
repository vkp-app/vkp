package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	authv1 "k8s.io/api/authorization/v1"
	"os"
)

func HasUser(ctx context.Context, _ any, next graphql.Resolver) (res any, err error) {
	log := logr.FromContextOrDiscard(ctx)
	_, ok := userctx.CtxUser(ctx)
	if !ok {
		log.V(1).Info("rejecting unauthorised request")
		return nil, ErrUnauthorised
	}
	return next(ctx)
}

func (r *Resolver) HasAdmin(ctx context.Context, _ any, next graphql.Resolver) (res any, err error) {
	log := logr.FromContextOrDiscard(ctx)
	user, ok := userctx.CtxUser(ctx)
	if !ok {
		log.V(1).Info("rejecting unauthorised request")
		return nil, ErrUnauthorised
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
		return nil, err
	}
	if !sar.Status.Allowed {
		log.Info("rejecting admin request due to failing SAR checks")
		return nil, ErrForbidden
	}
	return next(ctx)
}
