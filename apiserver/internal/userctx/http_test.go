package userctx_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Run("header is detected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://example.org", nil)
		req.Header.Set("X-Forwarded-User", "joe.bloggs")

		w := httptest.NewRecorder()

		userctx.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := userctx.CtxUser(r.Context())
			assert.True(t, ok)
			assert.EqualValues(t, "joe.bloggs", user.Username)
		})).ServeHTTP(w, req)
	})
	t.Run("no header is detected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://example.org", nil)
		w := httptest.NewRecorder()

		userctx.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := userctx.CtxUser(r.Context())
			assert.False(t, ok)
			assert.Nil(t, user)
		})).ServeHTTP(w, req)
	})
}
