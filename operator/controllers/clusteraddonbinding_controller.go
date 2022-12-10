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
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterAddonBindingReconciler reconciles a ClusterAddonBinding object
type ClusterAddonBindingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddonbindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddonbindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusteraddonbindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ClusterAddonBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("binding", req)
	log.Info("reconciling ClusterAddonBinding")

	br := &paasv1alpha1.ClusterAddonBinding{}
	if err := r.Get(ctx, req.NamespacedName, br); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve binding resource")
		return ctrl.Result{}, err
	}
	if br.DeletionTimestamp != nil {
		log.Info("skipping binding that has been marked for deletion")
		return ctrl.Result{}, nil
	}

	// reconcile
	if err := r.reconcileLabels(ctx, br); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileLabels copies the cluster and addon name
// into labels so that the API can find them using a label selector.
//
// It's a bit of a hack :/
func (r *ClusterAddonBindingReconciler) reconcileLabels(ctx context.Context, br *paasv1alpha1.ClusterAddonBinding) error {
	log := logging.FromContext(ctx)
	labels := br.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	var changed bool
	if val := labels[paasv1alpha1.LabelClusterRef]; val != br.Spec.ClusterRef.Name {
		labels[paasv1alpha1.LabelClusterRef] = br.Spec.ClusterRef.Name
		changed = true
	}
	if val := labels[paasv1alpha1.LabelClusterAddonRef]; val != br.Spec.ClusterAddonRef.Name {
		labels[paasv1alpha1.LabelClusterAddonRef] = br.Spec.ClusterAddonRef.Name
		changed = true
	}
	if changed {
		br.SetLabels(labels)
		if err := r.Update(ctx, br); err != nil {
			log.Error(err, "failed to update binding labels")
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterAddonBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.ClusterAddonBinding{}).
		Complete(r)
}
