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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PlatformSpec defines the desired state of Platform
type PlatformSpec struct {
	//+kubebuilder:validation:Required
	Namespace string `json:"namespace"`
	//+kubebuilder:validation:Required
	PrometheusURL string `json:"prometheusURL"`
	//+kubebuilder:validation:Required
	Domain string `json:"domain"`

	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy  corev1.PullPolicy             `json:"imagePullPolicy,omitempty"`

	Ingress ComponentIngressSpec `json:"ingress"`

	ApiServer ApiServerSpec `json:"apiServer,omitempty"`
	Dex       DexSpec       `json:"dex"`
}

type ApiServerSpec struct {
	ComponentSpec    `json:",inline"`
	OauthProxy       ComponentSpec `json:"oauthProxy,omitempty"`
	PrometheusConfig string        `json:"prometheusConfig,omitempty"`
}

type DexSpec struct {
	ComponentSpec `json:",inline"`
	Ingress       ComponentIngressSpec `json:"ingress"`
}

type ComponentSpec struct {
	Image           string                      `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty"`
	ReplicaCount    int32                       `json:"replicaCount,omitempty"`
	ExtraArgs       []string                    `json:"extraArgs,omitempty"`
}

type ComponentIngressSpec struct {
	Annotations      map[string]string           `json:"annotations,omitempty"`
	SecretRef        corev1.LocalObjectReference `json:"secretRef"`
	IngressClassName string                      `json:"ingressClassName,omitempty"`
}

// PlatformStatus defines the observed state of Platform
type PlatformStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Platform is the Schema for the platforms API
type Platform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformSpec   `json:"spec,omitempty"`
	Status PlatformStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PlatformList contains a list of Platform
type PlatformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Platform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Platform{}, &PlatformList{})
}
