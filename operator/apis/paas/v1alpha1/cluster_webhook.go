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
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/clusterutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var clusterlog = logf.Log.WithName("cluster-resource")

func (r *Cluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-cluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusters,verbs=create;update,versions=v1alpha1,name=mcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Cluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Cluster) Default() {
	clusterlog.Info("default", "name", r.Name)

	// generate the cluster ID
	if r.Status.ClusterID == "" {
		r.Status.ClusterID = clusterutil.NewID()
	}
	if r.Status.ClusterDomain == "" {
		r.Status.ClusterDomain = os.Getenv(clusterutil.EnvHostname)
	}
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-cluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=clusters,verbs=create;update,versions=v1alpha1,name=vcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Cluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateCreate() error {
	clusterlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateUpdate(old runtime.Object) error {
	clusterlog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	or := old.(*Cluster)

	// block the user from switching between HA and non-HA
	if or.Spec.HA.Enabled != r.Spec.HA.Enabled {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("ha").Child("enabled"), r.Spec.HA.Enabled, "high-availability cannot be changed"))
	}
	// block the user from changing the storage class
	if or.Spec.Storage.StorageClassName != r.Spec.Storage.StorageClassName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("storage").Child("storageClassName"), r.Spec.Storage.StorageClassName, "storage class cannot be changed"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindCluster}, r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateDelete() error {
	clusterlog.Info("validate delete", "name", r.Name)

	return nil
}
