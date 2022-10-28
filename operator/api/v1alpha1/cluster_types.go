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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	HA      HighAvailability `json:"ha,omitempty"`
	Storage Storage          `json:"storage,omitempty"`
}

type HighAvailability struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Storage struct {
	StorageClassName string `json:"storageClassName,omitempty"`
	Size             int    `json:"size,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	KubeVersion     string `json:"kubeVersion,omitempty"`
	PlatformVersion string `json:"platformVersion,omitempty"`

	KubeURL string `json:"kubeURL,omitempty"`
	WebURL  string `json:"webURL,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
