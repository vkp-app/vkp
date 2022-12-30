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

package paas

import (
	"context"
	"github.com/robfig/cron"
	"gitlab.dcas.dev/k8s/kube-glass/operator/internal/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AppliedClusterVersionReconciler reconciles a AppliedClusterVersion object
type AppliedClusterVersionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=appliedclusterversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=appliedclusterversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=appliedclusterversions/finalizers,verbs=update
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AppliedClusterVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling AppliedClusterVersion")

	// fetch the applied cluster version
	acv := &paasv1alpha1.AppliedClusterVersion{}
	if err := r.Get(ctx, req.NamespacedName, acv); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve AppliedClusterVersion")
		return ctrl.Result{}, err
	}
	if acv.DeletionTimestamp != nil {
		log.Info("skipping cluster that has been marked for deletion")
		return ctrl.Result{}, nil
	}

	// fetch the cluster
	cr := &paasv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: acv.Spec.ClusterRef.Name}, cr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve cluster resource")
		return ctrl.Result{}, err
	}

	// fetch the latest version
	cv, err := release.GetLatest(ctx, r.Client, cr.Spec.Track)
	if err != nil {
		return ctrl.Result{}, err
	}

	// check if we're within a maintenance window
	window, err := cron.ParseStandard(acv.Spec.MaintenanceWindow)
	if err != nil {
		log.Error(err, "failed to parse maintenance window CRON", "window", acv.Spec.MaintenanceWindow)
		return ctrl.Result{}, err
	}

	if !inWindow(window, time.Now()) {
		log.Info("skipping version update as we are not inside the maintenance window", "window", acv.Spec.MaintenanceWindow)
		r.Recorder.Eventf(acv, corev1.EventTypeNormal, "OutOfWindow", `Skipping maintenance check as we are outside the requested maintenance window. Next window: "%s"`, window.Next(time.Now()))
		return ctrl.Result{RequeueAfter: time.Hour}, nil
	}

	// return early if there have been no changes
	if acv.Status.VersionRef.Name == cv.GetName() {
		return ctrl.Result{}, nil
	}
	// update the status of our resource
	r.Recorder.Eventf(acv, corev1.EventTypeNormal, "Updated", `Updated cluster version (was: "%s", now: "%s")`, acv.Status.VersionRef.Name, cv.GetName())
	acv.Status.VersionRef.Name = cv.GetName()

	// update the status
	if err := r.Status().Update(ctx, acv); err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// inWindow returns whether a given time is within
// a CRON schedule.
func inWindow(s cron.Schedule, t time.Time) bool {
	next := s.Next(t)
	return next.Before(t.Add(time.Second * 90))
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppliedClusterVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.AppliedClusterVersion{}).
		Complete(r)
}
