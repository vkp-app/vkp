package userctx

import (
	"context"
	"github.com/go-logr/logr"
	"net/http"
	"strings"
)

type contextKey int

const (
	KeyUser   contextKey = iota
	KeyGroups contextKey = iota
)

const (
	headerUser   = "x-forwarded-user"
	headerGroups = "x-forwarded-groups"
)

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

// CtxUser extracts the users information from a
// given context.Context. If no user is present,
// a false value will be returned.
func CtxUser(ctx context.Context) (*User, bool) {
	username, ok := ctx.Value(KeyUser).(string)
	if !ok || username == "" {
		return nil, false
	}
	groups, ok := ctx.Value(KeyGroups).(string)
	if !ok {
		return nil, false
	}
	return &User{
		Username: username,
		// todo confirm that groups are comma-delimited
		Groups: strings.Split(groups, ","),
	}, true
}
