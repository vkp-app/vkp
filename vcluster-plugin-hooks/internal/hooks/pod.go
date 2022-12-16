package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const labelMetricTarget = "paas.dcas.dev/metric-target"

type PodHook struct {
	ClusterName string
}

func (h *PodHook) Name() string {
	return "pod-hook"
}

func (h *PodHook) Resource() client.Object {
	return &corev1.Pod{}
}

func (h *PodHook) MutateCreatePhysical(_ context.Context, obj client.Object) (client.Object, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object %+v is not a Pod", obj)
	}
	// add our labels
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Labels[labelMetricTarget] = h.ClusterName
	pod.Annotations[labelMetricTarget] = h.ClusterName
	return pod, nil
}

func (h *PodHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	return h.MutateCreatePhysical(ctx, obj)
}

var _ hook.ClientHook = &PodHook{}
var _ hook.MutateCreatePhysical = &PodHook{}
var _ hook.MutateUpdatePhysical = &PodHook{}
