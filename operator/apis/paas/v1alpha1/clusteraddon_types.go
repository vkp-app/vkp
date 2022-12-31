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

type AddonSource string

//goland:noinspection GoUnusedConst
const (
	SourceOfficial  AddonSource = "Official"
	SourcePlatform  AddonSource = "Platform"
	SourceCommunity AddonSource = "Community"
	SourceUnknown   AddonSource = "Unknown"
)

// ClusterAddonSpec defines the desired state of ClusterAddon
type ClusterAddonSpec struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resources"
	//+kubebuilder:validation:MinItems:=1
	Resources []RemoteRef `json:"resources"`

	// DisplayName is the human-readable name of the addon shown in
	// the addon marketplace.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Display name",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	DisplayName string `json:"displayName"`
	// Maintainer is the name/contact information of the addon.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Maintainer",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Maintainer string `json:"maintainer"`

	// Logo is the URL of an image that should be
	// shown in the addon marketplace for this addon.
	Logo string `json:"logo,omitempty"`
	// Description is the human-readable description of the addon
	// shown in the addon marketplace.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Description",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Description string `json:"description,omitempty"`
	// Source indicates where the addon came from and how
	// trustworthy it should be considered.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Source"
	//+kubebuilder:validation:Enum=Official;Platform;Community;Unknown
	Source AddonSource `json:"source,omitempty"`
	// SourceURL is an external HTTP address that can be used by users
	// to find more information about an addon.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Source URL",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	SourceURL string `json:"sourceURL,omitempty"`
}

type RemoteRef struct {
	// URL is a Kustomize-compatible HTTPS URL to a Kustomize directory. Mutually-exclusive with ConfigMap and Secret.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="URL",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	URL string `json:"url,omitempty"`
	// ConfigMap is a v1.ConfigMap that contains a number of Kustomize files. Mutually-exclusive with URL and Secret.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ConfigMap",xDescriptors="urn:alm:descriptor:io.kubernetes:ConfigMap"
	ConfigMap corev1.LocalObjectReference `json:"configMap,omitempty"`
	// Secret is a v1.Secret that contains a number of Kustomize files. Mutually-exclusive with URL and ConfigMap.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Secret",xDescriptors="urn:alm:descriptor:io.kubernetes:Secret"
	Secret corev1.LocalObjectReference `json:"secret,omitempty"`
	// OCI is an OCI-compliant container image that contains a Kustomize directory.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="OCI",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	OCI OCIRemoteRef `json:"oci,omitempty"`
}

type OCIRemoteRef struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image name",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Name string `json:"name,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image Pull Secret",xDescriptors="urn:alm:descriptor:io.kubernetes:Secret"
	ImagePullSecret corev1.LocalObjectReference `json:"imagePullSecret,omitempty"`
}

// ClusterAddonStatus defines the observed state of ClusterAddon
type ClusterAddonStatus struct {
	ResourceDigests map[string]string  `json:"resourceDigests,omitempty"`
	Conditions      []metav1.Condition `json:"conditions,omitempty"`
}

func ConfigMapDigestKey(name string) string {
	return "cm:" + name
}

func SecretDigestKey(name string) string {
	return "cm:" + name
}

func UriDigestKey(name string) string {
	return "uri:" + name
}

func OciDigestKey(name string) string {
	return "oci:" + name
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterAddon is the Schema for the clusteraddons API
type ClusterAddon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterAddonSpec   `json:"spec,omitempty"`
	Status ClusterAddonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterAddonList contains a list of ClusterAddon
type ClusterAddonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAddon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterAddon{}, &ClusterAddonList{})
}
