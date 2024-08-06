package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	// Name of the news source
	Name string `json:"name"`
	// Link to the news source
	Link string `json:"link"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	// Conditions represent the latest available observations of an object's state
	Status string `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Feed is the Schema for the feeds API
type Feed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeedSpec   `json:"spec,omitempty"`
	Status FeedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FeedList contains a list of Feed
type FeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Feed `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
