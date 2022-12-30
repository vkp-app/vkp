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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var clusteraddonbindinglog = logf.Log.WithName("clusteraddonbinding-resource")

func (r *ClusterAddonBinding) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-clusteraddonbinding,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusteraddonbindings,verbs=create;update,versions=v1alpha1,name=mclusteraddonbinding.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ClusterAddonBinding{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ClusterAddonBinding) Default() {
	clusteraddonbindinglog.Info("default", "name", r.Name)

	r.Labels[LabelClusterRef] = r.Spec.ClusterRef.Name
	r.Labels[LabelClusterAddonRef] = r.Spec.ClusterAddonRef.Name
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-clusteraddonbinding,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusteraddonbindings,verbs=create;update,versions=v1alpha1,name=vclusteraddonbinding.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ClusterAddonBinding{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddonBinding) ValidateCreate() error {
	clusteraddonbindinglog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList

	// block the user from changing the cluster ref
	if r.Spec.ClusterRef.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("clusterRef").Child("name"), "cluster reference must be set"))
	}
	// block the user from changing the addon ref
	if r.Spec.ClusterAddonRef.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("clusterAddonRef").Child("name"), "addon reference must be set"))
	}

	// validate that the name is in the appropriate format
	if r.Name != fmt.Sprintf("%s-%s", r.Spec.ClusterRef.Name, r.Spec.ClusterAddonRef.Name) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "resource name must follow the clustername-addonname format"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindClusterAddonBinding}, r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddonBinding) ValidateUpdate(old runtime.Object) error {
	clusteraddonbindinglog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	or := old.(*ClusterAddonBinding)

	// block the user from changing the cluster ref
	if or.Spec.ClusterRef.Name != r.Spec.ClusterRef.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("clusterRef").Child("name"), r.Spec.ClusterRef.Name, "cluster reference cannot be changed"))
	}
	// block the user from changing the addon ref
	if or.Spec.ClusterAddonRef.Name != r.Spec.ClusterAddonRef.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("clusterAddonRef").Child("name"), r.Spec.ClusterRef.Name, "addon reference cannot be changed"))
	}

	if r.Labels[LabelClusterRef] != r.Spec.ClusterRef.Name {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("metadata").Child("labels").Key(LabelClusterRef), "system fields cannot be changed"))
	}
	if r.Labels[LabelClusterAddonRef] != r.Spec.ClusterRef.Name {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("metadata").Child("labels").Key(LabelClusterAddonRef), "system fields cannot be changed"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindClusterAddonBinding}, r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddonBinding) ValidateDelete() error {
	clusteraddonbindinglog.Info("validate delete", "name", r.Name)

	return nil
}
