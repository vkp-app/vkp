package crdresources

import (
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/transutils"
	"k8s.io/apimachinery/pkg/api/equality"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceMonitorSyncer struct {
	translator.NamespacedTranslator
}

func NewServiceMonitorSyncer(ctx *synccontext.RegisterContext) *ServiceMonitorSyncer {
	return &ServiceMonitorSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, promv1.ServiceMonitorKindKey, &promv1.ServiceMonitor{}),
	}
}

func (s *ServiceMonitorSyncer) Init(ctx *synccontext.RegisterContext) error {
	return translate.EnsureCRDFromFile(ctx.Context, ctx.VirtualManager.GetConfig(), filepath.Join(os.Getenv(KoDataPathEnv), CustomResourceServiceMonitorFile), promv1.SchemeGroupVersion.WithKind(promv1.ServiceMonitorsKind))
}

func (s *ServiceMonitorSyncer) SyncDown(ctx *synccontext.SyncContext, vo client.Object) (ctrl.Result, error) {
	return s.SyncDownCreate(ctx, vo, s.translate(vo.(*promv1.ServiceMonitor)))
}

func (s *ServiceMonitorSyncer) Sync(ctx *synccontext.SyncContext, po, vo client.Object) (ctrl.Result, error) {
	vsm := vo.(*promv1.ServiceMonitor)
	psm := po.(*promv1.ServiceMonitor)

	return s.SyncDownUpdate(ctx, vo, s.translateUpdate(psm, vsm))
}

func (s *ServiceMonitorSyncer) translate(v client.Object) *promv1.ServiceMonitor {
	p := s.TranslateMetadata(v).(*promv1.ServiceMonitor)
	vsm := v.(*promv1.ServiceMonitor)
	p.Spec = *s.rewriteSpec(&vsm.Spec, vsm.Namespace)
	return p
}

func (s *ServiceMonitorSyncer) translateUpdate(po, vo *promv1.ServiceMonitor) *promv1.ServiceMonitor {
	var updated *promv1.ServiceMonitor

	// check annotations and labels
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(vo, po)
	if changed {
		updated = newIfNil(updated, po)
		updated.Annotations = updatedAnnotations
		updated.Labels = updatedLabels
	}

	// check spec
	ps := s.rewriteSpec(&vo.Spec, vo.GetNamespace())
	if !equality.Semantic.DeepEqual(*ps, po.Spec) {
		updated = newIfNil(updated, po)
		updated.Spec = *ps
	}

	return updated
}

func (s *ServiceMonitorSyncer) rewriteSpec(v *promv1.ServiceMonitorSpec, namespace string) *promv1.ServiceMonitorSpec {
	// translate selectors
	v = v.DeepCopy()
	if v.Selector.MatchLabels != nil {
		v.Selector.MatchLabels = transutils.TranslateLabels(v.Selector.MatchLabels, namespace)
	}
	for i, e := range v.Selector.MatchExpressions {
		v.Selector.MatchExpressions[i].Key = translator.ConvertLabelKey(e.Key)
	}
	if v.JobLabel != "" {
		v.JobLabel = translator.ConvertLabelKey(v.JobLabel)
	}
	return v
}

func newIfNil(updated *promv1.ServiceMonitor, po *promv1.ServiceMonitor) *promv1.ServiceMonitor {
	if updated == nil {
		return po.DeepCopy()
	}
	return updated
}

var _ syncer.Initializer = &ServiceMonitorSyncer{}
