/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TKMCacheClusterSpec defines the desired state of TKMCacheCluster
type TKMCacheClusterSpec struct {
	Name             string           `json:"name"`
	Image            string           `json:"image"`
	ResolvedDigest   string           `json:"resolvedDigest,omitempty"` // Injected by webhook or controller
	KernelProperties KernelProperties `json:"kernelProperties,omitempty"`
}

// TKMCacheClusterStatus defines the observed state of TKMCacheCluster
type TKMCacheClusterStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	LastSynced metav1.Time        `json:"lastSynced,omitempty"`
	Summary    []KernelSummary    `json:"summary,omitempty"`
	Digest     string             `json:"digest,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:webhook:path=/mutate-tkmcachecluster,mutating=true,failurePolicy=fail,groups=tkm.io,resources=tkmcacheclusters,verbs=create;update,versions=v1alpha1,name=mtkmcachecluster.kb.io,sideEffects=None,admissionReviewVersions=v1

// TKMCacheCluster is the Schema for the TKMCacheCluster API
type TKMCacheCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TKMCacheClusterSpec   `json:"spec,omitempty"`
	Status TKMCacheClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TKMCacheClusterList contains a list of TKMCacheCluster
type TKMCacheClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TKMCacheCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TKMCacheCluster{}, &TKMCacheClusterList{})
}
