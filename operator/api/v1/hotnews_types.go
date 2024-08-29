package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

// HotNewsSpec defines the desired state of HotNews.
type HotNewsSpec struct {
	Keywords      []string      `json:"keywords"`
	DateStart     string        `json:"dateStart,omitempty"`
	DateEnd       string        `json:"dateEnd,omitempty"`
	Feeds         []string      `json:"feeds,omitempty"`
	FeedGroups    []string      `json:"feedGroups,omitempty"`
	SummaryConfig SummaryConfig `json:"summaryConfig,omitempty"`
}

// SummaryConfig defines the configuration for summarizing news articles.
type SummaryConfig struct {
	TitlesCount int `json:"titlesCount,omitempty"`
}

// HotNewsStatus defines the observed state of HotNews.
type HotNewsStatus struct {
	ArticlesCount  int              `json:"articlesCount,omitempty"`
	NewsLink       string           `json:"newsLink,omitempty"`
	ArticlesTitles []string         `json:"articlesTitles,omitempty"`
	Condition      HotNewsCondition `json:"condition,omitempty"`
}

// HotNewsCondition represents the state of a Feed at a certain point.
type HotNewsCondition struct {
	// Status of the condition, one of True, False.
	Status bool `json:"status"`
	// If status is False, the reason should be populated
	Reason string `json:"reason,omitempty"`
}

func SetHotNewsErrorStatus(errorMessage string) HotNewsStatus {
	return HotNewsStatus{Condition: HotNewsCondition{
		Status: false,
		Reason: errorMessage,
	}}
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews represents resource and includes its specification and status.
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// ExtractFeedsFromGroups retrieves the feed names associated
// with the specified FeedGroups by referencing a ConfigMap.
func (h *HotNews) ExtractFeedsFromGroups(configMap v1.ConfigMap) []string {
	var feedNames []string

	for _, feedGroup := range h.Spec.FeedGroups {
		log.Printf("Processing FeedGroup: %s", feedGroup)
		if value, ok := configMap.Data[feedGroup]; ok {
			feedNames = append(feedNames, strings.Split(value, ",")...)
			log.Printf("Matched FeedGroup '%s' in ConfigMap, added values: %v", feedGroup, feedNames)
		}
	}

	return feedNames
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews resources.
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
