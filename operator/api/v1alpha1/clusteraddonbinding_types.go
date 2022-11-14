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

const (
	LabelClusterRef      = "paas.dcas.dev/cluster-name"
	LabelClusterAddonRef = "paas.dcas.dev/cluster-addon-name"
)

type AddonPhase string

const (
	AddonPhaseInstalling AddonPhase = "Installing"
	AddonPhaseInstalled  AddonPhase = "Installed"
	AddonPhaseDeleting   AddonPhase = "Deleting"
)

// ClusterAddonBindingSpec defines the desired state of ClusterAddonBinding
type ClusterAddonBindingSpec struct {
	ClusterRef      corev1.LocalObjectReference `json:"clusterRef"`
	ClusterAddonRef corev1.LocalObjectReference `json:"clusterAddonRef"`
}

// ClusterAddonBindingStatus defines the observed state of ClusterAddonBinding
type ClusterAddonBindingStatus struct {
	Phase AddonPhase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterAddonBinding is the Schema for the clusteraddonbindings API
type ClusterAddonBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterAddonBindingSpec   `json:"spec,omitempty"`
	Status ClusterAddonBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterAddonBindingList contains a list of ClusterAddonBinding
type ClusterAddonBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAddonBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterAddonBinding{}, &ClusterAddonBindingList{})
}
