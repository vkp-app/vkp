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
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/tenant"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Options TenantOptions
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants/finalizers,verbs=update
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("tenant", req.Name, "namespace", req.Namespace)
	log.Info("reconciling Tenant")

	tr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, req.NamespacedName, tr); err != nil {
		if errors.IsNotFound(err); err != nil {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}
	if tr.DeletionTimestamp != nil {
		log.Info("skipping tenant that has been marked for deletion")
		return ctrl.Result{}, nil
	}
	// basic reconciliation
	if err := r.reconcileCustomCA(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}
	if res, err := r.reconcileNamespaces(ctx, tr); err != nil || res.Requeue {
		return res, err
	}
	if err := r.reconcileAddons(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}

	// collect a list of managed clusters
	clusters := &paasv1alpha1.ClusterList{}
	var ns string
	// if the tenant uses a single namespace, we can limit
	// our search to just that namespace
	if tr.Spec.NamespaceStrategy == paasv1alpha1.StrategySingle {
		ns = tr.GetNamespace()
	}
	if tr.Status.Phase == "" {
		tr.Status.Phase = paasv1alpha1.PhasePendingApproval
	}
	if err := r.List(ctx, clusters, client.MatchingLabels{labelTenant: tr.GetName()}, client.InNamespace(ns)); err != nil {
		if errors.IsNotFound(err); err != nil {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to list managed clusters")
		return ctrl.Result{}, err
	}
	// generate the status object
	observedClusters := make([]paasv1alpha1.NamespacedName, len(clusters.Items))
	for i := range clusters.Items {
		observedClusters[i] = paasv1alpha1.NamespacedName{
			Namespace: clusters.Items[i].GetNamespace(),
			Name:      clusters.Items[i].GetName(),
		}
	}
	tr.Status.ObservedClusters = observedClusters
	// apply the update
	if err := r.Status().Update(ctx, tr); err != nil {
		log.Error(err, "failed to update tenant status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *TenantReconciler) reconcileCustomCA(ctx context.Context, tr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx).WithValues("tenant", tr.GetName())
	log.Info("reconciling CA")

	sec := tenant.CASecret(tr, r.Options.CustomCAFile)

	found := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tr.GetName(), Name: sec.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(tr, sec, r.Scheme)
			if err := r.Create(ctx, sec); err != nil {
				log.Error(err, "failed to create CA secret")
				return err
			}
			return nil
		}
		return err
	}
	if found.Data == nil {
		found.Data = sec.Data
		return r.Update(ctx, found)
	}
	if _, ok := found.Data[tenant.SecretKeyCA]; !ok {
		found.Data[tenant.SecretKeyCA] = sec.Data[tenant.SecretKeyCA]
		if err := r.Update(ctx, found); err != nil {
			log.Error(err, "failed to update CA secret")
			return err
		}
	}
	return nil
}

func (r *TenantReconciler) reconcileAddons(ctx context.Context, tr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx).WithValues("tenant", tr.GetName())
	log.Info("reconciling addons")

	if r.Options.SkipDefaultAddons {
		log.V(2).Info("skipping default addons as requested")
		return nil
	}

	found := cluster.Addons(tr)
	for _, addon := range found {
		if err := r.reconcileAddon(ctx, &addon, tr); err != nil {
			return err
		}
	}
	return nil
}

func (r *TenantReconciler) reconcileAddon(ctx context.Context, car *paasv1alpha1.ClusterAddon, tr *paasv1alpha1.Tenant) error {
	log := logging.FromContext(ctx).WithValues("addon", car.GetName(), "tenant", tr.GetName())
	log.Info("reconciling addon")

	found := &paasv1alpha1.ClusterAddon{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: tr.GetName(), Name: car.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(tr, car, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, car); err != nil {
				log.Error(err, "failed to create addon")
				return err
			}
			return nil
		}
		return err
	}
	if err := ctrl.SetControllerReference(tr, car, r.Scheme); err != nil {
		log.Error(err, "failed to set controller reference")
		return err
	}

	// reconcile any changes
	if !reflect.DeepEqual(car.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, car)
	}
	return nil
}

func (r *TenantReconciler) reconcileNamespaces(ctx context.Context, tr *paasv1alpha1.Tenant) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling namespaces")

	if tr.Spec.NamespaceStrategy == "" {
		tr.Spec.NamespaceStrategy = paasv1alpha1.StrategySingle
		if err := r.Update(ctx, tr); err != nil {
			log.Error(err, "failed to set tenant default namespace strategy")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// single namespaces have been pre-prepared (since they contain the tenant)
	// so we don't need to do anything
	if tr.Spec.NamespaceStrategy == paasv1alpha1.StrategySingle {
		log.Info("skipping namespace reconciliation due to namespace strategy")
		if len(tr.Status.ObservedNamespaces) == 0 || tr.Status.ObservedNamespaces[0] != tr.GetName() {
			tr.Status.ObservedNamespaces = []string{tr.GetName()}
			if err := r.Status().Update(ctx, tr); err != nil {
				log.Error(err, "failed to update tenant namespace list")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.Tenant{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *TenantReconciler) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.Update(ctx, new, option...)
}
