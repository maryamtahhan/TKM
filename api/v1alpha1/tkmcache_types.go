package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TKMCacheSpec defines the desired state of TKMCache
type TKMCacheSpec struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	ResolvedDigest string `json:"resolvedDigest,omitempty"` // Injected by webhook
}

type KernelProperties struct {
	TritonVersion string          `json:"tritonVersion"`
	Variant       string          `json:"variant,omitempty"`
	EntryCount    int             `json:"entryCount,omitempty"`
	Summary       []KernelSummary `json:"summary,omitempty"`
}

type KernelSummary struct {
	Backend  string `json:"backend"`
	Arch     string `json:"arch"`
	WarpSize int    `json:"warp_size"`
}

// TKMCacheStatus defines the observed state of TKMCache
type TKMCacheStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	LastSynced metav1.Time        `json:"lastSynced,omitempty"`
	Summary    []KernelSummary    `json:"summary,omitempty"`
	Digest     string             `json:"digest,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:webhook:path=/mutate-tkmcache,mutating=true,failurePolicy=fail,groups=tkm.io,resources=tkmcaches,verbs=create;update,versions=v1alpha1,name=mtkmcache.kb.io,sideEffects=None,admissionReviewVersions=v1

// TKMCache is the Schema for the TKMCache API
type TKMCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TKMCacheSpec   `json:"spec,omitempty"`
	Status TKMCacheStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TKMCacheList contains a list of TKMCache
type TKMCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TKMCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TKMCache{}, &TKMCacheList{})
}
