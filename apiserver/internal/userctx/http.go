package userctx

import (
	"context"
	"github.com/go-logr/logr"
	"net/http"
)

type contextKey int

const (
	KeyUser   contextKey = iota
	KeyGroups contextKey = iota
)

const headerUser = "x-forwarded-user"
const headerGroups = "x-forwarded-groups"

func Middleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logr.FromContextOrDiscard(r.Context())
			username := r.Header.Get(headerUser)
			groups := r.Header.Get(headerGroups)
			log = log.WithValues("username", username, "groups", groups)
			log.V(2).Info("extracting user middleware")
			// create the new context with information about
			// the user
			ctx := context.WithValue(logr.NewContext(r.Context(), log), KeyUser, username)
			ctx = context.WithValue(ctx, KeyGroups, groups)
			// continue as normal
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
