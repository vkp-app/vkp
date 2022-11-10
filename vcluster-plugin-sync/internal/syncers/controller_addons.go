package syncers

import (
	"context"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
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

func NewAddonSyncer(clusterName string) *AddonSyncer {
	return &AddonSyncer{
		ClusterName: clusterName,
	}
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
	for _, manifest := range cr.Spec.Resources {
		log.Info("applying manifest", "ManifestURL", manifest)
		var kustomizePath string
		// figure out where we get our resources from
		if manifest.URL != "" {
			// skip non-remote URLs
			if !strings.HasPrefix(manifest.URL, "https://") {
				log.Info("skipping non-HTTPS manifest")
				continue
			}
			kustomizePath = manifest.URL
		} else if manifest.ConfigMap.Name != "" {
			dir, err := r.manifestsFromConfigMap(ctx, manifest.ConfigMap.Name)
			if err != nil {
				return err
			}
			kustomizePath = dir
		} else if manifest.Secret.Name != "" {
			dir, err := r.manifestsFromSecret(ctx, manifest.Secret.Name)
			if err != nil {
				return err
			}
			kustomizePath = dir
		}
		// configure kustomize so we can
		// inflate helm charts
		opt := krusty.MakeDefaultOptions()
		opt.PluginConfig = ktypes.MakePluginConfig(ktypes.PluginRestrictionsNone, ktypes.BploLoadFromFileSys)
		opt.PluginConfig.FnpLoadingOptions.EnableStar = true
		opt.PluginConfig.HelmConfig.Enabled = true
		opt.PluginConfig.HelmConfig.Command = "helm"
		// render the manifests using kustomize
		allResources, err := krusty.MakeKustomizer(opt).Run(filesys.MakeFsOnDisk(), kustomizePath)
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
			obj, _, err := decoder.Decode([]byte(r.mangleYAML(string(yaml))), nil, nil)
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

// manifestsFromSecret downloads all kubernetes manifests from a corev1.ConfigMap
// and stores them in a temporary directory so that they can be consumed by
// Kustomize.
func (r *AddonSyncer) manifestsFromConfigMap(ctx context.Context, name string) (string, error) {
	log := logging.FromContext(ctx).WithValues("configmap", name)
	log.Info("fetching addon manifests from ConfigMap")

	cm := &corev1.ConfigMap{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: r.namespace, Name: name}, cm); err != nil {
		log.Error(err, "failed to retrieve ConfigMap")
		return "", err
	}
	dir, err := os.MkdirTemp("", "addon-*")
	if err != nil {
		log.Error(err, "failed to allocate temporary directory")
		return "", err
	}
	// write all the data to our temp directory
	for k, v := range cm.Data {
		if err := os.WriteFile(filepath.Join(dir, k), []byte(v), 0644); err != nil {
			log.Error(err, "failed to write file")
			return "", err
		}
	}
	return dir, nil
}

// manifestsFromSecret downloads all kubernetes manifests from a corev1.Secret
// and stores them in a temporary directory so that they can be consumed by
// Kustomize.
//
// Almost duplicate of manifestsFromConfigMap and could no doubt
// be improved down-the-line.
func (r *AddonSyncer) manifestsFromSecret(ctx context.Context, name string) (string, error) {
	log := logging.FromContext(ctx).WithValues("secret", name)
	log.Info("fetching addon manifests from Secret")

	sec := &corev1.Secret{}
	if err := r.pClient.Get(ctx, types.NamespacedName{Namespace: r.namespace, Name: name}, sec); err != nil {
		log.Error(err, "failed to retrieve Secret")
		return "", err
	}
	dir, err := os.MkdirTemp("", "addon-*")
	if err != nil {
		log.Error(err, "failed to allocate temporary directory")
		return "", err
	}
	// write all the data to our temp directory
	for k, v := range sec.Data {
		if err := os.WriteFile(filepath.Join(dir, k), v, 0644); err != nil {
			log.Error(err, "failed to write file")
			return "", err
		}
	}
	return dir, nil
}

// mangleYAML does a simple find-and-replace, so we can inject
// per-cluster configuration values (e.g. URLs and OAuth information)
func (*AddonSyncer) mangleYAML(s string) string {
	s = strings.ReplaceAll(s, MagicDexURL, os.Getenv(MagicDexURL))
	s = strings.ReplaceAll(s, MagicClusterURL, os.Getenv(MagicClusterURL))
	s = strings.ReplaceAll(s, MagicDexClientID, os.Getenv(MagicDexClientID))
	s = strings.ReplaceAll(s, MagicDexClientSecret, os.Getenv(MagicDexClientSecret))

	return s
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