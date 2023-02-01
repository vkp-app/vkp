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
	"strings"
)

// log is for logging in this package.
var tenantlog = logf.Log.WithName("tenant-resource")

func (r *Tenant) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-paas-dcas-dev-v1alpha1-tenant,mutating=true,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=tenants,verbs=create;update,versions=v1alpha1,name=mtenant.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Tenant{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Tenant) Default() {
	tenantlog.Info("default", "name", r.Name)

	if r.Spec.NamespaceStrategy == "" {
		r.Spec.NamespaceStrategy = StrategySingle
	}
}

//+kubebuilder:webhook:path=/validate-paas-dcas-dev-v1alpha1-tenant,mutating=false,failurePolicy=fail,sideEffects=None,groups=paas.dcas.dev,resources=tenants,verbs=create;update,versions=v1alpha1,name=vtenant.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Tenant{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateCreate() error {
	tenantlog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList

	if strings.HasPrefix(r.ObjectMeta.Name, "kube-") || strings.HasPrefix(r.ObjectMeta.Name, "openshift-") {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("name"), r.ObjectMeta.Name, "tenant name must not match kube-* or openshift-*"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindTenant}, r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateUpdate(old runtime.Object) error {
	tenantlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Tenant) ValidateDelete() error {
	tenantlog.Info("validate delete", "name", r.Name)

	return nil
}
