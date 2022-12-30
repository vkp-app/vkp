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

package v1

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
var oauthclientlog = logf.Log.WithName("oauthclient-resource")

func (r *OAuthClient) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-idp-dcas-dev-v1-oauthclient,mutating=true,failurePolicy=fail,sideEffects=None,groups=idp.dcas.dev,resources=oauthclients,verbs=create;update,versions=v1,name=moauthclient.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OAuthClient{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OAuthClient) Default() {
	oauthclientlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-idp-dcas-dev-v1-oauthclient,mutating=false,failurePolicy=fail,sideEffects=None,groups=idp.dcas.dev,resources=oauthclients,verbs=create;update,versions=v1,name=voauthclient.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &OAuthClient{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OAuthClient) ValidateCreate() error {
	oauthclientlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OAuthClient) ValidateUpdate(old runtime.Object) error {
	oauthclientlog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	or := old.(*OAuthClient)

	// block changing the client ID
	if or.Spec.ClientID != r.Spec.ClientID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec").Child("id"), "client id cannot be changed"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: KindOAuthClient}, r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OAuthClient) ValidateDelete() error {
	oauthclientlog.Info("validate delete", "name", r.Name)

	return nil
}
