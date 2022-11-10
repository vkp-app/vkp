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
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/kustomize/api/krusty"
	ktypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"strings"
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

	// get the cluster
	cr := &paasv1alpha1.Cluster{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, cr); err != nil {
		log.Error(err, "failed to retrieve cluster resource")
		return ctrl.Result{}, err
	}

	// reconcile the role binding
	if err := r.reconcileRBAC(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileAddons(ctx, cr); err != nil {
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

func (r *RoleBindingSyncer) reconcileAddons(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling nested addons")

	for _, addon := range cr.Spec.Addons {
		if err := r.applyAddon(ctx, types.NamespacedName{Namespace: addon.Namespace, Name: addon.Name}); err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleBindingSyncer) applyAddon(ctx context.Context, res types.NamespacedName) error {
	log := logging.FromContext(ctx).WithValues("namespace", res.Namespace, "name", res.Name)
	log.Info("applying nested addon")

	// fetch the addon
	car := &paasv1alpha1.ClusterAddon{}
	if err := r.pClient.Get(ctx, res, car); err != nil {
		log.Error(err, "failed to retrieve addon resource")
		return err
	}
	// apply it
	for _, manifest := range car.Spec.Manifests {
		log.Info("applying manifest", "ManifestURL", manifest)
		// skip non-remote URLs
		if !strings.HasPrefix(manifest, "https://") {
			log.Info("skipping non-HTTPS manifest")
			continue
		}
		// apply the manifests using kustomize
		allResources, err := krusty.MakeKustomizer(&krusty.Options{
			AddManagedbyLabel: true,
			DoPrune:           true,
			PluginConfig: &ktypes.PluginConfig{
				HelmConfig: ktypes.HelmConfig{
					Enabled: true,
				},
			},
		}).Run(filesys.MakeFsInMemory(), manifest)
		if err != nil {
			log.Error(err, "failed to kustomize addon")
			return err
		}
		for _, resource := range allResources.Resources() {
			log = log.WithValues("resourceName", resource.GetName(), "resourceKind", resource.GetKind(), "resourceNamespace", resource.GetNamespace())
			log.Info("applying resource")
			yaml, err := resource.AsYAML()
			if err != nil {
				log.Error(err, "failed to convert resource to YAML")
				return err
			}
			decoder := clientgoscheme.Codecs.UniversalDeserializer()
			obj, _, err := decoder.Decode(yaml, nil, nil)
			if err != nil {
				log.Error(err, "failed to decode resource YAML")
				return err
			}
			if err := r.vClient.Update(ctx, obj.(client.Object)); err != nil {
				log.Error(err, "failed to apply resource")
				return err
			}
		}
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
