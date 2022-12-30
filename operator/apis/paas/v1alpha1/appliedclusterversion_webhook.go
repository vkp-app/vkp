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

package v1alpha1

import (
	"github.com/robfig/cron"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var appliedclusterversionlog = logf.Log.WithName("appliedclusterversion-resource")

func (r *AppliedClusterVersion) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-appliedclusterversion,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=appliedclusterversions,verbs=create;update,versions=v1alpha1,name=mappliedclusterversion.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &AppliedClusterVersion{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *AppliedClusterVersion) Default() {
	appliedclusterversionlog.Info("default", "name", r.Name)

	if r.Spec.MaintenanceWindow == "" {
		r.Spec.MaintenanceWindow = "* * * * *"
	}
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-appliedclusterversion,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=appliedclusterversions,verbs=create;update,versions=v1alpha1,name=vappliedclusterversion.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AppliedClusterVersion{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AppliedClusterVersion) ValidateCreate() error {
	appliedclusterversionlog.Info("validate create", "name", r.Name)

	return r.validateAppliedClusterVersion()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AppliedClusterVersion) ValidateUpdate(runtime.Object) error {
	appliedclusterversionlog.Info("validate update", "name", r.Name)

	return r.validateAppliedClusterVersion()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AppliedClusterVersion) ValidateDelete() error {
	appliedclusterversionlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *AppliedClusterVersion) validateAppliedClusterVersion() error {
	var allErrs field.ErrorList
	// validate the cron schedule
	if _, err := cron.ParseStandard(r.Spec.MaintenanceWindow); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("maintenanceWindow"), r.Spec.MaintenanceWindow, err.Error()))
	}
	if r.Spec.ClusterRef.Name != r.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("clusterRef").Child("name"), r.Spec.ClusterRef.Name, "resource name and cluster name must be identical"))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindAppliedClusterVersion}, r.Name, allErrs)
}
