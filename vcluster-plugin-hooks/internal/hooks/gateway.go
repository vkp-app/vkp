package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"strings"
)

type GatewayHook struct {
	ClusterDomain string
}

func (h *GatewayHook) Name() string {
	return "gateway-hook"
}

func (h *GatewayHook) Resource() client.Object {
	return &gatewayv1beta1.Gateway{}
}

func (h *GatewayHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	log := logging.FromContext(ctx)
	gateway, ok := obj.(*gatewayv1beta1.Gateway)
	if !ok {
		return nil, fmt.Errorf("object %+v is not a Gateway", obj)
	}
	// collect all the hosts
	var r string
	for i, l := range gateway.Spec.Listeners {
		// skip gateways that don't specify
		// a hostname
		if l.Hostname == nil {
			continue
		}
		r = string(*l.Hostname)
		if strings.HasSuffix(r, "."+h.ClusterDomain) || r == h.ClusterDomain {
			continue
		}
		flat := strings.ReplaceAll(r, ".", "-")
		hostname := gatewayv1beta1.Hostname(fmt.Sprintf("%s.%s", flat, h.ClusterDomain))
		gateway.Spec.Listeners[i].Hostname = &hostname
		log.Info("rewriting gateway.networking.k8s.io/Gateway hostname", "old", r, "new", hostname)
	}
	return gateway, nil
}

func (h *GatewayHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	return h.MutateCreatePhysical(ctx, obj)
}

var _ hook.ClientHook = &GatewayHook{}
var _ hook.MutateCreatePhysical = &GatewayHook{}
var _ hook.MutateUpdatePhysical = &GatewayHook{}
