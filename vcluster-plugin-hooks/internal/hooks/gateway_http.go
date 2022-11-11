package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"strings"
)

type GatewayHttpHook struct {
	ClusterDomain string
}

func (h *GatewayHttpHook) Name() string {
	return "gateway-http-hook"
}

func (h *GatewayHttpHook) Resource() client.Object {
	return &gatewayv1beta1.HTTPRoute{}
}

func (h *GatewayHttpHook) MutateCreatePhysical(_ context.Context, obj client.Object) (client.Object, error) {
	route, ok := obj.(*gatewayv1beta1.HTTPRoute)
	if !ok {
		return nil, fmt.Errorf("object %+v is not an HTTPRoute", obj)
	}
	// collect all the hosts
	var r string
	for i := range route.Spec.Hostnames {
		r = string(route.Spec.Hostnames[i])
		if strings.HasSuffix(r, "."+h.ClusterDomain) {
			continue
		}
		flat := strings.ReplaceAll(r, ".", "-")
		route.Spec.Hostnames[i] = gatewayv1beta1.Hostname(fmt.Sprintf("%s.%s", flat, h.ClusterDomain))
	}
	return route, nil
}

func (h *GatewayHttpHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	return h.MutateCreatePhysical(ctx, obj)
}

var _ hook.ClientHook = &GatewayHttpHook{}
var _ hook.MutateCreatePhysical = &GatewayHttpHook{}
var _ hook.MutateUpdatePhysical = &GatewayHttpHook{}
