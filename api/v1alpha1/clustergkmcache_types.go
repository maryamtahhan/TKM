/*
Copyright 2025.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterGKMCacheSpec defines the desired state of ClusterGKMCache
type ClusterGKMCacheSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ClusterGKMCache. Edit clustergkmcache_types.go to remove/update
	Name             string           `json:"name"`
	Image            string           `json:"image"`
	ResolvedDigest   string           `json:"resolvedDigest,omitempty"` // Injected by webhook or controller
	KernelProperties KernelProperties `json:"kernelProperties,omitempty"`
}

// ClusterGKMCacheStatus defines the observed state of ClusterGKMCache
type ClusterGKMCacheStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	LastSynced metav1.Time        `json:"lastSynced,omitempty"`
	Summary    []KernelSummary    `json:"summary,omitempty"`
	Digest     string             `json:"digest,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterGKMCache is the Schema for the clustergkmcaches API
type ClusterGKMCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterGKMCacheSpec   `json:"spec,omitempty"`
	Status ClusterGKMCacheStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterGKMCacheList contains a list of ClusterGKMCache
type ClusterGKMCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterGKMCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterGKMCache{}, &ClusterGKMCacheList{})
}
