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
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"reflect"
	capiv1betav1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vclusters,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("cluster", req.NamespacedName)
	log.Info("reconciling Cluster")

	cr := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve cluster resource")
		return ctrl.Result{}, err
	}
	if cr.DeletionTimestamp != nil {
		log.Info("skipping cluster that has been marked for deletion")
		return ctrl.Result{}, nil
	}

	if err := r.reconcileID(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileDomain(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileVCluster(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileCluster(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileIngress(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.Cluster{}).
		Owns(&vclusterv1alpha1.VCluster{}).
		Owns(&capiv1betav1.Cluster{}).
		Owns(&netv1.Ingress{}).
		Complete(r)
}

func (r *ClusterReconciler) reconcileID(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.V(1).Info("reconciling cluster ID")

	if cr.Status.ClusterID == "" {
		clusterID := cluster.NewID()
		log.Info("generated cluster ID", "ID", clusterID)
		cr.Status.ClusterID = clusterID
		return r.Status().Update(ctx, cr)
	}
	return nil
}

func (r *ClusterReconciler) reconcileDomain(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.V(1).Info("reconciling cluster domain")

	// set the cluster domain
	// so that we can manage it independently
	// of the operator
	if cr.Status.ClusterDomain == "" {
		cr.Status.ClusterDomain = os.Getenv(cluster.EnvHostname)
		return r.Status().Update(ctx, cr)
	}
	return nil
}

func (r *ClusterReconciler) reconcileVCluster(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling vcluster")

	vcluster, err := cluster.VCluster(ctx, cr)
	if err != nil {
		return err
	}

	found := &vclusterv1alpha1.VCluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, vcluster, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, vcluster); err != nil {
				log.Error(err, "failed to create vcluster")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(vcluster.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, vcluster)
	}
	return nil
}

func (r *ClusterReconciler) reconcileCluster(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling cluster")

	capiCluster := cluster.Cluster(cr)

	found := &capiv1betav1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, capiCluster, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, capiCluster); err != nil {
				log.Error(err, "failed to create cluster")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(capiCluster.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, capiCluster)
	}
	return nil
}

func (r *ClusterReconciler) reconcileIngress(ctx context.Context, cr *paasv1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling ingress")

	ing := cluster.Ingress(cr)

	found := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, ing, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, ing); err != nil {
				log.Error(err, "failed to create ingress")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(ing.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, ing)
	}
	return nil
}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *ClusterReconciler) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.Update(ctx, new, option...)
}
