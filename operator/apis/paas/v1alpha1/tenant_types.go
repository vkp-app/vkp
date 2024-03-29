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

type NamespaceStrategy string

//goland:noinspection GoUnusedConst
const (
	StrategySingle   NamespaceStrategy = "Single"
	StrategyMultiple NamespaceStrategy = "Multiple"
)

type TenantPhase string

const (
	PhasePendingApproval TenantPhase = "PendingApproval"
	PhaseReady           TenantPhase = "Ready"
)

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	// Owner is the human user that owns the tenant.
	// They will have special privileges that will not be given
	// to Accessors (e.g. ability to delete the tenant).
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Owner",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Owner             string            `json:"owner"`
	NamespaceStrategy NamespaceStrategy `json:"namespaceStrategy"`
	// Accessors define who is authorised to interact with the tenant.
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Accessors"
	Accessors []AccessRef `json:"accessors,omitempty"`
}

// TenantStatus defines the observed state of Tenant
type TenantStatus struct {
	ObservedClusters   []NamespacedName   `json:"observedClusters,omitempty"`
	ObservedNamespaces []string           `json:"observedNamespaces,omitempty"`
	Phase              TenantPhase        `json:"phase,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

type NamespacedName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Tenant is the Schema for the tenants API
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.owner`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
