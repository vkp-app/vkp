package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-logr/logr"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
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
