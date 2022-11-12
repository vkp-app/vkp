package hooks

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"testing"
)

func TestGatewayHttpHook_MutateCreatePhysical(t *testing.T) {
	h := &GatewayHttpHook{ClusterDomain: "example.org"}
	hostname := gatewayv1beta1.Hostname("foo.example.io")
	hostname2 := gatewayv1beta1.Hostname("foo.example.org")

	t.Run("invalid object is rejected", func(t *testing.T) {
		_, err := h.MutateCreatePhysical(context.TODO(), &corev1.Pod{})
		assert.Error(t, err)
	})
	t.Run("hostnames are rewritten", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &gatewayv1beta1.HTTPRoute{
			Spec: gatewayv1beta1.HTTPRouteSpec{
				Hostnames: []gatewayv1beta1.Hostname{
					hostname,
					hostname2,
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "foo-example-io.example.org", ing.(*gatewayv1beta1.HTTPRoute).Spec.Hostnames[0])
		assert.EqualValues(t, "foo.example.org", ing.(*gatewayv1beta1.HTTPRoute).Spec.Hostnames[1])
	})
}
