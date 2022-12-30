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

// AppliedClusterVersionSpec defines the desired state of AppliedClusterVersion
type AppliedClusterVersionSpec struct {
	ClusterRef corev1.LocalObjectReference `json:"clusterRef"`
	//+kubebuilder:default:="* * * * *"
	MaintenanceWindow string `json:"maintenanceWindow,omitempty"`
}

// AppliedClusterVersionStatus defines the observed state of AppliedClusterVersion
type AppliedClusterVersionStatus struct {
	VersionRef corev1.LocalObjectReference `json:"versionRef,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AppliedClusterVersion is the Schema for the appliedclusterversions API
type AppliedClusterVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppliedClusterVersionSpec   `json:"spec,omitempty"`
	Status AppliedClusterVersionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AppliedClusterVersionList contains a list of AppliedClusterVersion
type AppliedClusterVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppliedClusterVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppliedClusterVersion{}, &AppliedClusterVersionList{})
}
