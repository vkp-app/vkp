package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
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

func (h *IngressHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	log := logging.FromContext(ctx)
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
		log.Info("rewriting networking.k8s.io/Ingress hostname", "old", r.Host, "new", ing.Spec.Rules[i].Host)
	}
	// collect all the TLS hosts
	for i, t := range ing.Spec.TLS {
		for j, tt := range t.Hosts {
			if strings.HasSuffix(tt, "."+h.ClusterDomain) {
				continue
			}
			flat := strings.ReplaceAll(tt, ".", "-")
			ing.Spec.TLS[i].Hosts[j] = fmt.Sprintf("%s.%s", flat, h.ClusterDomain)
			log.Info("rewriting networking.k8s.io/Ingress TLS hostname", "old", tt, "new", ing.Spec.TLS[i].Hosts[j])
		}
	}
	return ing, nil
}

func (h *IngressHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	return h.MutateCreatePhysical(ctx, obj)
}

var _ hook.ClientHook = &IngressHook{}
var _ hook.MutateCreatePhysical = &IngressHook{}
var _ hook.MutateUpdatePhysical = &IngressHook{}
