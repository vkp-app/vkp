package hooks

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"testing"
)

func TestIngressHook_MutateCreatePhysical(t *testing.T) {
	h := &IngressHook{ClusterDomain: "example.org"}

	t.Run("invalid object is rejected", func(t *testing.T) {
		_, err := h.MutateCreatePhysical(context.TODO(), &corev1.Pod{})
		assert.Error(t, err)
	})
	t.Run("hostnames are rewritten", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &netv1.Ingress{
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "foo.domain.io",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "foo-domain-io.example.org", ing.(*netv1.Ingress).Spec.Rules[0].Host)
	})
	t.Run("tls hostnames are rewritten", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &netv1.Ingress{
			Spec: netv1.IngressSpec{
				TLS: []netv1.IngressTLS{
					{
						Hosts: []string{
							"foo.domain.io",
						},
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "foo-domain-io.example.org", ing.(*netv1.Ingress).Spec.TLS[0].Hosts[0])
	})
	t.Run("prepared hostnames are skipped", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &netv1.Ingress{
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "foo.example.org",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "foo.example.org", ing.(*netv1.Ingress).Spec.Rules[0].Host)
	})
	t.Run("prepared tls hostnames are skipped", func(t *testing.T) {
		ing, err := h.MutateCreatePhysical(context.TODO(), &netv1.Ingress{
			Spec: netv1.IngressSpec{
				TLS: []netv1.IngressTLS{
					{
						Hosts: []string{
							"foo.example.org",
						},
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.EqualValues(t, "foo.example.org", ing.(*netv1.Ingress).Spec.TLS[0].Hosts[0])
	})
}
