package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	// Name of the news source
	Name string `json:"name"`
	// Link of the news sources
	Link string `json:"link"`
}

// ConditionType represents a condition type for a Feed
type ConditionType string

const (
	// ConditionAdded indicates that the feed has been successfully added
	ConditionAdded ConditionType = "Added"
	// ConditionUpdated indicates that the feed has been successfully updated
	ConditionUpdated ConditionType = "Updated"
	// ConditionDeleted indicates that the feed has been successfully deleted
	ConditionDeleted ConditionType = "Deleted"
)

// Condition represents the state of a Feed at a certain point.
type Condition struct {
	// Type of the condition, e.g., Added, Updated, Deleted.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False.
	Status bool `json:"status"`
	// If status is False, the reason should be populated
	Reason string `json:"reason,omitempty"`
	// If status is False, the message should be populated
	Message string `json:"message,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// AddCondition adds a new condition to the FeedStatus and updates all conditions' LastUpdateTime.
func (f *FeedStatus) AddCondition(condition Condition) {
	f.Conditions = append(f.Conditions, condition)
	f.updateConditions()
}

// Contains checks if a condition of the specified type with the given status
// exists in the FeedStatus conditions
func (f *FeedStatus) Contains(conditionType ConditionType, status bool) bool {
	for _, condition := range f.Conditions {
		if condition.Type == conditionType && condition.Status == status {
			return true
		}
	}
	return false
}

// updateConditions updates the LastUpdateTime for all conditions in FeedStatus.
func (f *FeedStatus) updateConditions() {
	for i := range f.Conditions {
		f.Conditions[i].LastUpdateTime = metav1.Now()
	}
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []Condition `json:"conditions,omitempty"`
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
