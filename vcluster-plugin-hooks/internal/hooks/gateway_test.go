package hooks

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"testing"
)

func TestGatewayHook_MutateCreatePhysical(t *testing.T) {
	h := &GatewayHook{ClusterDomain: "example.org"}
	hostname := gatewayv1beta1.Hostname("example.io")
	hostname2 := gatewayv1beta1.Hostname(h.ClusterDomain)

	t.Run("invalid object is rejected", func(t *testing.T) {
		_, err := h.MutateCreatePhysical(context.TODO(), &corev1.Pod{})
		assert.Error(t, err)
	})
	t.Run("hostnames are rewritten", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &gatewayv1beta1.Gateway{
			Spec: gatewayv1beta1.GatewaySpec{
				Listeners: []gatewayv1beta1.Listener{
					{
						Name:     "foobar",
						Hostname: &hostname,
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "example-io.example.org", *ing.(*gatewayv1beta1.Gateway).Spec.Listeners[0].Hostname)
	})
	t.Run("prepared hostnames are skipped", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &gatewayv1beta1.Gateway{
			Spec: gatewayv1beta1.GatewaySpec{
				Listeners: []gatewayv1beta1.Listener{
					{
						Name:     "foobar",
						Hostname: &hostname2,
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "example.org", *ing.(*gatewayv1beta1.Gateway).Spec.Listeners[0].Hostname)
	})
}
