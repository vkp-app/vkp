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
var clusterversionlog = logf.Log.WithName("clusterversion-resource")

func (r *ClusterVersion) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-clusterversion,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusterversions,verbs=create;update,versions=v1alpha1,name=mclusterversion.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ClusterVersion{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ClusterVersion) Default() {
	clusterversionlog.Info("default", "name", r.Name)
	r.Labels[LabelTrackRef] = string(r.Spec.Track)
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-clusterversion,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusterversions,verbs=create;update;delete,versions=v1alpha1,name=vclusterversion.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ClusterVersion{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterVersion) ValidateCreate() error {
	clusterversionlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterVersion) ValidateUpdate(old runtime.Object) error {
	clusterversionlog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	or := old.(*ClusterVersion)

	// block changes to the chart version
	if or.Spec.Chart.Version != r.Spec.Chart.Version {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec").Child("chart").Child("version"), "chart version cannot be changed - create a new cluster version instead"))
	}
	// block changes to the image tag.
	// we can change the repo/registry as we would
	// expect that they serve the same image
	if or.Spec.Image.Tag != r.Spec.Image.Tag {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec").Child("image").Child("tag"), "image tag cannot be changed - create a new cluster version instead"))
	}

	if r.Labels[LabelTrackRef] != string(r.Spec.Track) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("metadata").Child("labels").Key(LabelTrackRef), "system fields cannot be changed"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindCluster}, r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterVersion) ValidateDelete() error {
	clusterversionlog.Info("validate delete", "name", r.Name)

	// straight-up block deletion
	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindClusterVersion}, r.Name, []*field.Error{
		field.Forbidden(field.NewPath("spec"), "cluster versions cannot be deleted as it could destroy active clusters"),
	})
}
