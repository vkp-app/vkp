package syncers

import (
	"context"
	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/nested"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

type CertificateSyncer struct {
	pClient     client.Client
	vClient     client.Client
	Scheme      *runtime.Scheme
	ClusterName string
	namespace   string
}

const (
	namespaceTrust = "kube-system"
	secretTrust    = "tls-identity"
)

func NewCertificateSyncer(clusterName string) *CertificateSyncer {
	return &CertificateSyncer{
		ClusterName: clusterName,
	}
}

func (*CertificateSyncer) Name() string {
	return "certificate-syncer"
}

func (*CertificateSyncer) Resource() client.Object {
	return &certv1.Certificate{}
}

func (r *CertificateSyncer) Register(ctx *synccontext.RegisterContext) error {
	r.pClient = ctx.PhysicalManager.GetClient()
	r.vClient = ctx.VirtualManager.GetClient()
	r.Scheme = ctx.PhysicalManager.GetScheme()
	r.namespace = ctx.CurrentNamespace

	return ctrl.NewControllerManagedBy(ctx.PhysicalManager).
		For(&certv1.Certificate{}).
		Complete(r)
}

func (r *CertificateSyncer) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx, "namespace", req.Namespace, "name", req.Name)
	log.Info("reconciling certificate")

	// skip resources that don't match
	// our namespace
	if req.Namespace != r.namespace || req.Name != cluster.IdentitySecretName(r.ClusterName) {
		log.Info("skipping certificate outside our domain")
		return ctrl.Result{}, nil
	}

	// fetch the secret
	sec := &corev1.Secret{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: r.namespace, Name: cluster.IdentitySecretName(r.ClusterName)}, sec); err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("failed to locate trust secret as it has not been issues")
			return ctrl.Result{Requeue: true}, nil
		}
		log.Error(err, "failed to retrieve trust secret resource")
		return ctrl.Result{}, err
	}
	newSec := nested.TrustCertificate(sec.Data)

	found := &corev1.Secret{}
	if err := r.vClient.Get(ctx, types.NamespacedName{Namespace: newSec.GetNamespace(), Name: newSec.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.vClient.Create(ctx, newSec); err != nil {
				log.Error(err, "failed to create secret")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err := r.updateSecretData(ctx, "tls.crt", found, newSec); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.updateSecretData(ctx, "tls.key", found, newSec); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.updateSecretData(ctx, "ca.crt", found, newSec); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *CertificateSyncer) updateSecretData(ctx context.Context, key string, found, sec *corev1.Secret) error {
	log := logging.FromContext(ctx)
	existing := sec.Data[key]
	if val, ok := found.Data[key]; !ok || val == nil {
		found.Data[key] = existing
		if err := r.vClient.Update(ctx, found); err != nil {
			log.Error(err, "failed to update secret")
			return err
		}
	}
	return nil
}

var _ syncer.Base = &CertificateSyncer{}
var _ syncer.ControllerStarter = &CertificateSyncer{}
