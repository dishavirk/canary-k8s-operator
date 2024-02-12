/*
Copyright 2024.

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

// Package v1alpha1 defines the version of the API group.
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CanarySpec defines the desired state of Canary
// It specifies what we want
type CanarySpec struct {
	DeploymentName string `json:"deploymentName"`
	Image          string `json:"image"` // Add this line
	Percentage     int    `json:"percentage"`
	Replicas       int32  `json:"replicas"` // Ensure this is defined for specifying replica count
}

// CanaryStatus defines the observed state of Canary
// It reflects the current state of the canary deployment as observed by the operator
type CanaryStatus struct {
	Phase string   `json:"phase"`
	Nodes []string `json:"nodes,omitempty"` // Add this line
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Canary is the Schema for the canaries API
// This top-level type represents a Canary resource in K8s
type Canary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CanarySpec   `json:"spec,omitempty"`
	Status CanaryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CanaryList contains a list of Canary.
// This type is used to list or watch multiple Canary resources
type CanaryList struct {
	metav1.TypeMeta `json:",inline"`            // Includes API version and kind for the list
	metav1.ListMeta `json:"metadata,omitempty"` // Standard list metadata
	Items           []Canary                    `json:"items"` // The list of Canary resources
}

func init() {
	SchemeBuilder.Register(&Canary{}, &CanaryList{}) // Registers the Canary and CanaryList types with the runtime scheme
}
