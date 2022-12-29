package crdresources

import (
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OAuthClientSyncer struct {
	translator.NamespacedTranslator
}

func NewOAuthClientSyncer(ctx *synccontext.RegisterContext) *OAuthClientSyncer {
	return &OAuthClientSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, "oauthclient", &idpv1.OAuthClient{}),
	}
}

func (s *OAuthClientSyncer) Init(ctx *synccontext.RegisterContext) error {
	return translate.EnsureCRDFromFile(ctx.Context, ctx.VirtualManager.GetConfig(), filepath.Join(os.Getenv(KoDataPathEnv), CustomResourceOAuthClientFile), idpv1.GroupVersion.WithKind("OAuthClient"))
}

func (s *OAuthClientSyncer) SyncDown(ctx *synccontext.SyncContext, vo client.Object) (ctrl.Result, error) {
	// skip resources that the user has specifically
	// requested that we don't sync
	if vo.GetAnnotations()[LabelDisableSync] != "" || vo.GetLabels()[LabelDisableSync] != "" {
		return ctrl.Result{}, nil
	}
	return s.SyncDownCreate(ctx, vo, s.translate(vo.(*idpv1.OAuthClient)))
}

func (s *OAuthClientSyncer) Sync(ctx *synccontext.SyncContext, po, vo client.Object) (ctrl.Result, error) {
	// skip resources that the user has specifically
	// requested that we don't sync
	if vo.GetAnnotations()[LabelDisableSync] != "" || vo.GetLabels()[LabelDisableSync] != "" {
		return ctrl.Result{}, nil
	}

	vsm := vo.(*idpv1.OAuthClient)
	psm := po.(*idpv1.OAuthClient)

	return s.SyncDownUpdate(ctx, vo, s.translateUpdate(psm, vsm))
}

func (s *OAuthClientSyncer) translate(v client.Object) *idpv1.OAuthClient {
	p := s.TranslateMetadata(v).(*idpv1.OAuthClient)
	vsm := v.(*idpv1.OAuthClient)
	p.Spec = *s.rewriteSpec(&vsm.Spec, vsm.GetNamespace())
	return p
}

func (s *OAuthClientSyncer) translateUpdate(po, vo *idpv1.OAuthClient) *idpv1.OAuthClient {
	var updated *idpv1.OAuthClient

	// check annotations and labels
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(vo, po)
	if changed {
		updated = s.newIfNil(updated, po)
		updated.Annotations = updatedAnnotations
		updated.Labels = updatedLabels
	}

	// check spec
	ps := s.rewriteSpec(&vo.Spec, vo.GetNamespace())
	if !equality.Semantic.DeepEqual(*ps, po.Spec) {
		updated = s.newIfNil(updated, po)
		updated.Spec = *ps
	}

	return updated
}

func (s *OAuthClientSyncer) rewriteSpec(v *idpv1.OAuthClientSpec, namespace string) *idpv1.OAuthClientSpec {
	// translate selectors
	v = v.DeepCopy()
	v.ClientSecretRef.Name = translate.PhysicalName(v.ClientSecretRef.Name, namespace)
	return v
}

func (*OAuthClientSyncer) newIfNil(updated *idpv1.OAuthClient, po *idpv1.OAuthClient) *idpv1.OAuthClient {
	if updated == nil {
		return po.DeepCopy()
	}
	return updated
}

var _ syncer.Initializer = &OAuthClientSyncer{}
