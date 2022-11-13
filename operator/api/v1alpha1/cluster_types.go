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

type ReleaseTrack int

const (
	TrackStable  ReleaseTrack = iota
	TrackRegular ReleaseTrack = iota
	TrackRapid   ReleaseTrack = iota
	TrackBeta    ReleaseTrack = iota
)

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	HA      HighAvailability `json:"ha,omitempty"`
	Storage Storage          `json:"storage,omitempty"`

	// Track defines the frequency of system updates.
	Track ReleaseTrack `json:"track,omitempty"`

	// Accessors define who is authorised to interact with the cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Accessors"
	Accessors []AccessRef `json:"accessors,omitempty"`
}

type HighAvailability struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled",xDescriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	Enabled bool `json:"enabled,omitempty"`
}

type Storage struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage class",xDescriptors="urn:alm:descriptor:io.kubernetes:StorageClass"
	StorageClassName string `json:"storageClassName,omitempty"`
	// Size in Gi of the clusters backing disk.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size (Gi)",xDescriptors="urn:alm:descriptor:com.tectonic.ui:number"
	Size int `json:"size,omitempty"`
}

type AccessRef struct {
	// ReadOnly indicates that the user/group should only have view access
	// to the virtual cluster.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Read Only",xDescriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	ReadOnly bool `json:"readOnly,omitempty"`
	// User binds a user to the virtual cluster. Mutually-exclusive with, and preferred over Group.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="User",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	User string `json:"user,omitempty"`
	// Group binds a group to the virtual cluster.
	// Using groups should be preferred as it allows you to manage
	// membership outside of Kubernetes. Mutually-exclusive with User.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Group",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Group string `json:"group,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Cluster conditions",xDescriptors="urn:alm:descriptor:io.kubernetes.conditions"
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	KubeVersion     string `json:"kubeVersion,omitempty"`
	PlatformVersion string `json:"platformVersion,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Kubernetes API URL",xDescriptors="urn:alm:descriptor:org.w3:link"
	KubeURL string `json:"kubeURL,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dashboard URL",xDescriptors="urn:alm:descriptor:org.w3:link"
	WebURL string `json:"webURL,omitempty"`

	ClusterID     string `json:"clusterID,omitempty"`
	ClusterDomain string `json:"clusterDomain,omitempty"`

	Message StatusMessage `json:"message,omitempty"`

	Inventory NestedInventory `json:"inventory,omitempty"`
}

type NestedInventory struct {
	AccessorRoles []string `json:"accessorRoles,omitempty"`
}

type StatusMessage struct {
	Kube   string `json:"kube,omitempty"`
	Addons string `json:"addons,omitempty"`
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
