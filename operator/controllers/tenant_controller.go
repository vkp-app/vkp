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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=tenants/finalizers,verbs=update
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("tenant", req.Name, "namespace", req.Namespace)
	log.Info("reconciling Tenant")

	cr := &paasv1alpha1.Tenant{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		if errors.IsNotFound(err); err != nil {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve tenant resource")
		return ctrl.Result{}, err
	}
	if cr.DeletionTimestamp != nil {
		log.Info("skipping tenant that has been marked for deletion")
		return ctrl.Result{}, nil
	}
	// basic reconciliation
	if res, err := r.reconcileNamespaces(ctx, cr); err != nil {
		return ctrl.Result{}, err
	} else if res.Requeue {
		return res, nil
	}

	// collect a list of managed clusters
	clusters := &paasv1alpha1.ClusterList{}
	var ns string
	// if the tenant uses a single namespace, we can limit
	// our search to just that namespace
	if cr.Spec.NamespaceStrategy == paasv1alpha1.StrategySingle {
		ns = cr.GetNamespace()
	}
	if cr.Status.Phase == "" {
		cr.Status.Phase = paasv1alpha1.PhasePendingApproval
	}
	if err := r.List(ctx, clusters, client.MatchingLabels{labelTenant: cr.GetName()}, client.InNamespace(ns)); err != nil {
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
	cr.Status.ObservedClusters = observedClusters
	// apply the update
	if err := r.Status().Update(ctx, cr); err != nil {
		log.Error(err, "failed to update tenant status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *TenantReconciler) reconcileNamespaces(ctx context.Context, cr *paasv1alpha1.Tenant) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling namespaces")

	if cr.Spec.NamespaceStrategy == "" {
		cr.Spec.NamespaceStrategy = paasv1alpha1.StrategySingle
		if err := r.Update(ctx, cr); err != nil {
			log.Error(err, "failed to set tenant default namespace strategy")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// single namespaces have been pre-prepared (since they contain the tenant)
	// so we don't need to do anything
	if cr.Spec.NamespaceStrategy == paasv1alpha1.StrategySingle {
		log.Info("skipping namespace reconciliation due to namespace strategy")
		if len(cr.Status.ObservedNamespaces) == 0 || cr.Status.ObservedNamespaces[0] != cr.GetNamespace() {
			cr.Status.ObservedNamespaces = []string{cr.GetNamespace()}
			if err := r.Status().Update(ctx, cr); err != nil {
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
		Complete(r)
}
