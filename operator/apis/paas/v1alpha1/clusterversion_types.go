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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
)

const LabelTrackRef = "paas.dcas.dev/release-track"

// ClusterVersionSpec defines the desired state of ClusterVersion
type ClusterVersionSpec struct {
	//+kubebuilder:validation:Required
	Image ClusterVersionImage `json:"image"`
	Chart ClusterVersionChart `json:"chart,omitempty"`
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=Stable;Regular;Rapid;Beta
	Track ReleaseTrack `json:"track"`
}

type ClusterVersionChart struct {
	Repository string `json:"repository,omitempty"`
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
}

type ClusterVersionImage struct {
	Registry string `json:"registry,omitempty"`
	//+kubebuilder:validation:Required
	Repository string `json:"repository"`
	//+kubebuilder:validation:Required
	Tag string `json:"tag"`
}

func (in *ClusterVersionImage) String() string {
	base := fmt.Sprintf("%s:%s", in.Repository, in.Tag)
	if in.Registry == "" {
		return base
	}
	return filepath.Join(in.Registry, base)
}

// ClusterVersionStatus defines the observed state of ClusterVersion
type ClusterVersionStatus struct {
	VersionNumber ClusterVersionNumber `json:"versionNumber,omitempty"`
	Conditions    []metav1.Condition   `json:"conditions,omitempty"`
}

type ClusterVersionNumber struct {
	Major int64 `json:"major,omitempty"`
	Minor int64 `json:"minor,omitempty"`
	Patch int64 `json:"patch,omitempty"`
}

func (in *ClusterVersionNumber) String() string {
	return fmt.Sprintf("%d.%d", in.Major, in.Minor)
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ClusterVersion is the Schema for the clusterversions API
// +kubebuilder:printcolumn:name="Tag",type=string,JSONPath=`.spec.image.tag`
type ClusterVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterVersionSpec   `json:"spec,omitempty"`
	Status ClusterVersionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterVersionList contains a list of ClusterVersion
type ClusterVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterVersion{}, &ClusterVersionList{})
}
