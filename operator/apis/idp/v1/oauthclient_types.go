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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OAuthClientSpec defines the desired state of OAuthClient
type OAuthClientSpec struct {
	ClientID        string                   `json:"id"`
	ClientSecretRef corev1.SecretKeySelector `json:"secretRef"`

	RedirectURIs []string `json:"redirectURIs"`

	TrustedPeers []string `json:"trustedPeers,omitempty"`
	Public       bool     `json:"public,omitempty"`
	LogoURL      string   `json:"logoURL,omitempty"`
}

// OAuthClientStatus defines the observed state of OAuthClient
type OAuthClientStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OAuthClient is the Schema for the oauthclients API
type OAuthClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OAuthClientSpec   `json:"spec,omitempty"`
	Status OAuthClientStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OAuthClientList contains a list of OAuthClient
type OAuthClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OAuthClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OAuthClient{}, &OAuthClientList{})
}
