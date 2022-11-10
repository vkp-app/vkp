package syncers

import (
	"context"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/nested"
	authv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

type RoleBindingSyncer struct {
	pClient client.Client
	vClient client.Client
	Scheme  *runtime.Scheme
}

func NewRBACSyncer() *RoleBindingSyncer {
	return &RoleBindingSyncer{}
}

func (*RoleBindingSyncer) Name() string {
	return "role-binding-syncer"
}

func (*RoleBindingSyncer) Resource() client.Object {
	return &paasv1alpha1.Cluster{}
}

func (r *RoleBindingSyncer) Register(ctx *synccontext.RegisterContext) error {
	r.pClient = ctx.PhysicalManager.GetClient()
	r.vClient = ctx.VirtualManager.GetClient()
	r.Scheme = ctx.PhysicalManager.GetScheme()

	return ctrl.NewControllerManagedBy(ctx.PhysicalManager).
		For(&paasv1alpha1.Cluster{}).
		Complete(r)
}

func (r *RoleBindingSyncer) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx, "namespace", req.Namespace, "name", req.Name)
	log.Info("syncing up")

	// get the tenant
	tr := &paasv1alpha1.Tenant{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Namespace}, tr); err != nil {
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}

	// reconcile the role binding
	if err := r.reconcileRBAC(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *RoleBindingSyncer) reconcileRBAC(ctx context.Context, cr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested RBAC")

	binding := nested.OwnerBinding(cr)

	found := &authv1.ClusterRoleBinding{}
	if err := r.vClient.Get(ctx, types.NamespacedName{Name: binding.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.vClient.Create(ctx, binding); err != nil {
				log.Error(err, "failed to create nested ClusterRoleBinding")
				return err
			}
			return nil
		}
		log.Error(err, "failed to retrieve virtual role binding")
		return err
	}
	// reconcile changes
	if !reflect.DeepEqual(found.Subjects, binding.Subjects) {
		return r.SafeUpdate(ctx, found, binding)
	}
	return nil
}

var _ syncer.Base = &RoleBindingSyncer{}
var _ syncer.ControllerStarter = &RoleBindingSyncer{}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *RoleBindingSyncer) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.vClient.Update(ctx, new, option...)
}
