package syncers

import (
	"context"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/nested"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	cr := &paasv1alpha1.Tenant{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Namespace}, cr); err != nil {
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}

	// reconcile the role binding
	if err := r.reconcileRBAC(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *RoleBindingSyncer) reconcileRBAC(ctx context.Context, cr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested RBAC")

	binding := nested.OwnerBinding(cr)
	if err := r.vClient.Create(ctx, binding); err != nil {
		log.Error(err, "failed to create nested ClusterRoleBinding")
		return err
	}
	return nil
}

var _ syncer.Base = &RoleBindingSyncer{}
var _ syncer.ControllerStarter = &RoleBindingSyncer{}
