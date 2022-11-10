package syncers

import (
	"context"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/kustomize/api/krusty"
	ktypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"strings"
)

type AddonSyncer struct {
	pClient     client.Client
	vClient     client.Client
	Scheme      *runtime.Scheme
	ClusterName string
	namespace   string
}

func NewAddonSyncer() *AddonSyncer {
	return &AddonSyncer{}
}

func (*AddonSyncer) Name() string {
	return "addon-syncer"
}

func (*AddonSyncer) Resource() client.Object {
	return &paasv1alpha1.ClusterAddon{}
}

func (r *AddonSyncer) Register(ctx *synccontext.RegisterContext) error {
	r.pClient = ctx.PhysicalManager.GetClient()
	r.vClient = ctx.VirtualManager.GetClient()
	r.Scheme = ctx.PhysicalManager.GetScheme()
	r.namespace = ctx.CurrentNamespace

	return ctrl.NewControllerManagedBy(ctx.PhysicalManager).
		For(&paasv1alpha1.ClusterAddon{}).
		Complete(r)
}

func (r *AddonSyncer) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx, "namespace", req.Namespace, "name", req.Name)
	log.Info("reconciling addon")

	// get the cluster
	//cr := &paasv1alpha1.Cluster{}
	//if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, cr); err != nil {
	//	log.Error(err, "failed to retrieve cluster resource")
	//	return ctrl.Result{}, err
	//}

	// skip resources that don't match
	// our namespace
	if req.Namespace != r.namespace {
		log.Info("skipping addon outside our domain")
		return ctrl.Result{}, nil
	}

	// fetch the addon
	car := &paasv1alpha1.ClusterAddon{}
	if err := r.pClient.Get(ctx, req.NamespacedName, car); err != nil {
		log.Error(err, "failed to retrieve addon resource")
		return ctrl.Result{}, err
	}

	// reconcile the addon
	if err := r.reconcileAddon(ctx, car); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AddonSyncer) reconcileAddon(ctx context.Context, cr *paasv1alpha1.ClusterAddon) error {
	log := logging.FromContext(ctx).WithValues("namespace", cr.Namespace, "name", cr.Name)
	log.Info("reconciling nested addon")

	// apply it
	for _, manifest := range cr.Spec.Manifests {
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
		}).Run(filesys.MakeFsOnDisk(), manifest)
		if err != nil {
			log.Error(err, "failed to kustomize addon")
			return err
		}
		for _, resource := range allResources.Resources() {
			log = log.WithValues("resourceName", resource.GetName(), "resourceKind", resource.GetKind(), "resourceNamespace", resource.GetNamespace())
			log.Info("applying resource")
			// set the namespace if one is not set
			if resource.GetNamespace() == "" {
				_ = resource.SetNamespace("default")
			}
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
			// update or create the resource
			res := obj.(client.Object)
			if err := r.vClient.Update(ctx, res); err != nil {
				if errors.IsNotFound(err) {
					if err := r.vClient.Create(ctx, res); err != nil {
						log.Error(err, "failed to create resource")
						return err
					}
					continue
				}
				log.Error(err, "failed to update resource")
				return err
			}
		}
	}
	return nil
}

var _ syncer.Base = &AddonSyncer{}
var _ syncer.ControllerStarter = &AddonSyncer{}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *AddonSyncer) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.vClient.Update(ctx, new, option...)
}
