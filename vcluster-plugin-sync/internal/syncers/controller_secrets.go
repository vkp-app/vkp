package syncers

import (
	"context"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/tenant"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/nested"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

type SecretSyncer struct {
	clusterName string
	namespace   string
	pClient     client.Client
	vClient     client.Client
	Scheme      *runtime.Scheme
}

func NewSecretSyncer(clusterName, namespace string) *SecretSyncer {
	return &SecretSyncer{
		clusterName: clusterName,
		namespace:   namespace,
	}
}

func (*SecretSyncer) Name() string {
	return "secret-syncer"
}

func (*SecretSyncer) Resource() client.Object {
	return &corev1.Namespace{}
}

func (r *SecretSyncer) Register(ctx *synccontext.RegisterContext) error {
	r.pClient = ctx.PhysicalManager.GetClient()
	r.vClient = ctx.VirtualManager.GetClient()
	r.Scheme = ctx.PhysicalManager.GetScheme()

	return ctrl.NewControllerManagedBy(ctx.VirtualManager).
		For(&corev1.Namespace{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func (r *SecretSyncer) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx, "name", req.Name)
	log.Info("syncing secrets")

	// get the namespace
	ns := &corev1.Namespace{}
	if err := r.vClient.Get(ctx, req.NamespacedName, ns); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve namespace")
		return ctrl.Result{}, err
	}

	// get the tenant
	tr := &paasv1alpha1.Tenant{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Name: r.namespace}, tr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}

	// reconcile the role binding
	if err := r.reconcileRootCA(ctx, ns, tr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SecretSyncer) reconcileRootCA(ctx context.Context, namespace *corev1.Namespace, tr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling CA secret")

	// fetch the CA secret
	sec := &corev1.Secret{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: r.namespace, Name: tenant.CASecretName(tr.GetName())}, sec); err != nil {
		log.Error(err, "failed to retrieve physical CA secret")
		return err
	}

	cm := nested.RootCA(namespace.GetName(), string(sec.Data[tenant.SecretKeyCA]))

	found := &corev1.ConfigMap{}
	if err := r.vClient.Get(ctx, types.NamespacedName{Name: cm.GetName(), Namespace: namespace.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			log.Info("creating virtual root CA")
			_ = ctrl.SetControllerReference(namespace, cm, r.vClient.Scheme())
			if err := r.vClient.Create(ctx, cm); err != nil {
				log.Error(err, "failed to create nested ConfigMap")
				return err
			}
			return nil
		}
		log.Error(err, "failed to retrieve virtual ConfigMap")
		return err
	}
	_ = ctrl.SetControllerReference(namespace, cm, r.vClient.Scheme())
	data := found.Data[tenant.SecretKeyCA]
	newData := cm.Data[tenant.SecretKeyCA]
	if newData == "" {
		return nil
	}
	if data != newData {
		log.Info("patching root CA")
		return r.SafeUpdate(ctx, found, cm)
	}
	return nil
}

var _ syncer.Base = &SecretSyncer{}
var _ syncer.ControllerStarter = &SecretSyncer{}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *SecretSyncer) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.vClient.Update(ctx, new, option...)
}
