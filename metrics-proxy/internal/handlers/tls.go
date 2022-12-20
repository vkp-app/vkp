package handlers

import (
	"context"
	"github.com/go-logr/logr"
	"net/http"
	"strings"
)

type ContextKeyRequester int

const (
	RequestingCluster ContextKeyRequester = iota
	RequestingTenant  ContextKeyRequester = iota
)

func TLS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logr.FromContextOrDiscard(r.Context())

		commonName := r.TLS.VerifiedChains[0][0].Subject.CommonName
		tenant, cluster, _ := strings.Cut(commonName, "/")

		log.Info("extracted request information from TLS certificate", "tenant", tenant, "cluster", cluster)

		ctx := context.WithValue(r.Context(), RequestingTenant, tenant)
		ctx = context.WithValue(ctx, RequestingCluster, cluster)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
