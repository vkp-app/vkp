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
	"github.com/Masterminds/semver"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

// ClusterVersionReconciler reconciles a ClusterVersion object
type ClusterVersionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusterversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusterversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusterversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ClusterVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("clusterVersion", req.NamespacedName)
	log.Info("reconciling ClusterVersion")

	cvr := &paasv1alpha1.ClusterVersion{}
	if err := r.Get(ctx, req.NamespacedName, cvr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve cluster version resource")
		return ctrl.Result{}, err
	}
	if cvr.DeletionTimestamp != nil {
		log.Info("skipping cluster version that has been marked for deletion")
		return ctrl.Result{}, nil
	}

	// reconcile labels that we use
	// for listing
	if val := cvr.ObjectMeta.Labels[paasv1alpha1.LabelTrackRef]; val != string(cvr.Spec.Track) {
		// panic guard
		if cvr.ObjectMeta.Labels == nil {
			cvr.ObjectMeta.Labels = map[string]string{}
		}
		cvr.ObjectMeta.Labels[paasv1alpha1.LabelTrackRef] = string(cvr.Spec.Track)
		if err := r.Update(ctx, cvr); err != nil {
			log.Error(err, "failed to update cluster version")
			return ctrl.Result{}, err
		}
	}
	// parse the semantic version
	sv, err := semver.NewVersion(cvr.Spec.Image.Tag)
	if err != nil {
		log.Error(err, "failed to parse semantic version", "tag", cvr.Spec.Image.Tag)
		return ctrl.Result{}, err
	}

	// update the CR to include the semantic
	// versions we parsed out of it
	cvr.Status.VersionNumber = paasv1alpha1.ClusterVersionNumber{
		Major: sv.Major(),
		Minor: sv.Minor(),
		Patch: sv.Patch(),
	}
	if err := r.Status().Update(ctx, cvr); err != nil {
		log.Error(err, "failed to update cluster version status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.ClusterVersion{}).
		Complete(r)
}