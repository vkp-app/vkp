package syncers

import (
	"context"
	"github.com/djcass44/go-utils/utilities/sliceutils"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/nested"
	authv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// get the cluster
	cr := &paasv1alpha1.Cluster{}
	if err := r.pClient.Get(ctx, req.NamespacedName, cr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve cluster resource")
		return ctrl.Result{}, err
	}

	// get the tenant
	tr := &paasv1alpha1.Tenant{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Namespace}, tr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}

	// reconcile the role binding
	if err := r.reconcileOwnerBinding(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}
	if res, err := r.reconcileClusterAccessorBinding(ctx, cr); err != nil || res.Requeue {
		return res, err
	}

	return ctrl.Result{}, nil
}

func (r *RoleBindingSyncer) reconcileClusterAccessorBinding(ctx context.Context, cr *paasv1alpha1.Cluster) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested accessor RBAC")

	for _, ar := range cr.Spec.Accessors {
		if err := r.reconcileAccessor(ctx, cr, &ar); err != nil {
			log.Error(err, "failed to reconcile accessor")
			return ctrl.Result{}, err
		}
	}
	if err := r.pClient.Status().Update(ctx, cr); err != nil {
		log.Error(err, "failed to update cluster inventory status")
		return ctrl.Result{}, err
	}

	// check if there are any roles that we need to purge
	if len(cr.Spec.Accessors) == len(cr.Status.Inventory.AccessorRoles) {
		return ctrl.Result{}, nil
	}

	log.Info("detected mismatch in inventoried and actual accessors", "expected", len(cr.Spec.Accessors), "actual", len(cr.Status.Inventory.AccessorRoles))

	// delete all the roles and request a requeue
	for _, role := range cr.Status.Inventory.AccessorRoles {
		res := &authv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: role,
			},
		}
		log.V(1).Info("deleting cluster role binding")
		if err := r.vClient.Delete(ctx, res); err != nil && !errors.IsNotFound(err) {
			log.Error(err, "failed to delete cluster role binding")
			return ctrl.Result{}, err
		}
		// remove this role
		cr.Status.Inventory.AccessorRoles = sliceutils.Filter(cr.Status.Inventory.AccessorRoles, func(s string) bool {
			return s != role
		})
		if err := r.pClient.Status().Update(ctx, cr); err != nil {
			log.Error(err, "failed to update cluster inventory status")
			return ctrl.Result{}, err
		}
	}
	// request a requeue so we can
	// recreate the expected bindings
	return ctrl.Result{Requeue: true}, nil
}

func (r *RoleBindingSyncer) reconcileAccessor(ctx context.Context, cr *paasv1alpha1.Cluster, ar *paasv1alpha1.AccessRef) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested accessor", "Resource", ar)

	binding := nested.AccessWriteBinding(ar)

	if !sliceutils.Includes(cr.Status.Inventory.AccessorRoles, binding.GetName()) {
		cr.Status.Inventory.AccessorRoles = append(cr.Status.Inventory.AccessorRoles, binding.GetName())
	}

	// fetch the current resource
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
	if !reflect.DeepEqual(found.Subjects, binding.Subjects) || !reflect.DeepEqual(found.RoleRef, binding.RoleRef) {
		return r.SafeUpdate(ctx, found, binding)
	}
	return nil
}

func (r *RoleBindingSyncer) reconcileOwnerBinding(ctx context.Context, tr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested owner RBAC")

	binding := nested.OwnerBinding(tr)

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
