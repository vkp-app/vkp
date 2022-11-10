package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type IngressHook struct {
	ClusterDomain string
}

func (h *IngressHook) Name() string {
	return "ingress-hook"
}

func (h *IngressHook) Resource() client.Object {
	return &netv1.Ingress{}
}

func (h *IngressHook) MutateCreatePhysical(_ context.Context, obj client.Object) (client.Object, error) {
	ing, ok := obj.(*netv1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object %+v is not an Ingress", obj)
	}
	// collect all the hosts
	for i, r := range ing.Spec.Rules {
		if strings.HasSuffix(r.Host, "."+h.ClusterDomain) {
			continue
		}
		flat := strings.ReplaceAll(r.Host, ".", "-")
		ing.Spec.Rules[i].Host = fmt.Sprintf("%s.%s", flat, h.ClusterDomain)
	}
	return ing, nil
}

func (h *IngressHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	return h.MutateCreatePhysical(ctx, obj)
}

var _ hook.ClientHook = &IngressHook{}
var _ hook.MutateCreatePhysical = &IngressHook{}
var _ hook.MutateUpdatePhysical = &IngressHook{}
