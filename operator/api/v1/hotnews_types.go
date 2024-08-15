package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ArticlesCount  int      `json:"articlesCount,omitempty"`
	NewsLink       string   `json:"newsLink,omitempty"`
	ArticlesTitles []string `json:"articlesTitles,omitempty"`
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
