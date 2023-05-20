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
	"os"
	"regexp"
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

	// validate the name of the tenant.
	// see (#80)
	matchers := []nameMatchFunc{r.regexMatch, r.prefixMatch}
	for _, m := range matchers {
		err := m(r.ObjectMeta.Name)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindTenant}, r.Name, allErrs)
}

func (*Tenant) prefixMatch(name string) *field.Error {
	if strings.HasPrefix(name, "kube-") || strings.HasPrefix(name, "openshift-") {
		return field.Invalid(field.NewPath("metadata").Child("name"), name, "tenant name must not match kube-* or openshift-*")
	}
	return nil
}

func (*Tenant) regexMatch(name string) *field.Error {
	customRegex := os.Getenv("TENANT_NAME_REGEX")
	if customRegex == "" {
		return nil
	}
	tenantlog.V(1).Info("parsing tenant name regex", "expr", customRegex)
	expr, err := regexp.Compile(customRegex)
	if err != nil {
		tenantlog.Error(err, "unable to compile tenant name regex")
		return nil
	}
	// if the tenant name doesn't match the custom
	// regexp, reject it
	if !expr.MatchString(name) {
		return field.Invalid(field.NewPath("metadata").Child("name"), name, "tenant name does not match requirements")
	}
	return nil
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
