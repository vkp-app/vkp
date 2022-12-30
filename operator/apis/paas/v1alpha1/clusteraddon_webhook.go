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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var clusteraddonlog = logf.Log.WithName("clusteraddon-resource")

func (r *ClusterAddon) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-clusteraddon,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusteraddons,verbs=create;update,versions=v1alpha1,name=mclusteraddon.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ClusterAddon{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ClusterAddon) Default() {
	clusteraddonlog.Info("default", "name", r.Name)

	if r.Spec.Source == "" {
		r.Spec.Source = SourceUnknown
	}
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-clusteraddon,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusteraddons,verbs=create;update,versions=v1alpha1,name=vclusteraddon.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ClusterAddon{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddon) ValidateCreate() error {
	clusteraddonlog.Info("validate create", "name", r.Name)

	return r.validateClusterAddon()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddon) ValidateUpdate(old runtime.Object) error {
	clusteraddonlog.Info("validate update", "name", r.Name)

	return r.validateClusterAddon()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterAddon) ValidateDelete() error {
	clusteraddonlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *ClusterAddon) validateClusterAddon() error {
	var allErrs field.ErrorList

	// validate that resources comply with the OneOf policy
	for i, res := range r.Spec.Resources {
		var resourceCount int
		if res.URL != "" {
			resourceCount++
		}
		if res.ConfigMap.Name != "" {
			resourceCount++
		}
		if res.Secret.Name != "" {
			resourceCount++
		}
		if res.OCI.Name != "" {
			resourceCount++
		}
		if resourceCount == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("resources"), "a resource must be provided"))
		} else if resourceCount > 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("resources").Index(i), res, "resource must only contain one resource type"))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindClusterAddon}, r.Name, allErrs)
}
