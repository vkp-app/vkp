/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

// ClusterAddonReconciler reconciles a ClusterAddon object
type ClusterAddonReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ClusterAddonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("addon", req.NamespacedName)
	log.Info("reconciling ClusterAddon")

	car := &paasv1alpha1.ClusterAddon{}
	if err := r.Get(ctx, req.NamespacedName, car); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve ClusterAddon resource")
		return ctrl.Result{}, err
	}
	if car.DeletionTimestamp != nil {
		log.Info("skipping addon that has been marked for deletion")
		return ctrl.Result{}, nil
	}

	// generate digests for all resources
	car.Status.ResourceDigests = map[string]string{}

	for _, rr := range car.Spec.Resources {
		if rr.URL != "" {
			if err := r.reconcileUri(ctx, car, &rr); err != nil {
				log.Error(err, "failed to reconcile remote URL")
				continue
			}
		} else if rr.ConfigMap.Name != "" {
			if err := r.reconcileConfigMap(ctx, car, &rr); err != nil {
				log.Error(err, "failed to reconcile ConfigMap")
				continue
			}
		} else if rr.Secret.Name != "" {
			if err := r.reconcileSecret(ctx, car, &rr); err != nil {
				log.Error(err, "failed to reconcile Secret")
				continue
			}
		}
	}

	if len(car.Status.ResourceDigests) == 0 {
		return ctrl.Result{}, nil
	}

	// update the status
	if err := r.Status().Update(ctx, car); err != nil {
		log.Error(err, "failed to update addon status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ClusterAddonReconciler) reconcileConfigMap(ctx context.Context, car *paasv1alpha1.ClusterAddon, res *paasv1alpha1.RemoteRef) error {
	return r.reconcileResource(ctx, car, &corev1.ConfigMap{}, paasv1alpha1.ConfigMapDigestKey, res.ConfigMap.Name)
}

func (r *ClusterAddonReconciler) reconcileSecret(ctx context.Context, car *paasv1alpha1.ClusterAddon, res *paasv1alpha1.RemoteRef) error {
	return r.reconcileResource(ctx, car, &corev1.Secret{}, paasv1alpha1.SecretDigestKey, res.Secret.Name)
}

func (r *ClusterAddonReconciler) reconcileResource(ctx context.Context, car *paasv1alpha1.ClusterAddon, res client.Object, keygen func(string) string, name string) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling addon resource", "resourceKind", res.GetObjectKind().GroupVersionKind().String(), "resourceName", res.GetName())

	if err := r.Get(ctx, types.NamespacedName{Namespace: car.GetNamespace(), Name: name}, res); err != nil {
		return err
	}
	// take ownership so we can listen to
	// update events
	if err := controllerutil.SetOwnerReference(car, res, r.Scheme); err != nil {
		log.Error(err, "failed to take ownership of resource")
		return err
	}
	if err := r.Update(ctx, res); err != nil {
		log.Error(err, "failed to update resource")
		return err
	}
	// generate and store the resource hash
	car.Status.ResourceDigests[keygen(res.GetName())] = res.GetResourceVersion()
	return nil
}

func (r *ClusterAddonReconciler) reconcileUri(ctx context.Context, car *paasv1alpha1.ClusterAddon, res *paasv1alpha1.RemoteRef) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling addon URI")

	// generate and store the resource hash
	car.Status.ResourceDigests[paasv1alpha1.UriDigestKey(res.URL)] = res.URL
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterAddonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.ClusterAddon{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
