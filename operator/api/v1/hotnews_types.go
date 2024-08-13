package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HotNewsSpec struct {
	Keywords      []string      `json:"keywords"`
	DateStart     string        `json:"dateStart,omitempty"`
	DateEnd       string        `json:"dateEnd,omitempty"`
	Feeds         []string      `json:"feeds,omitempty"`
	FeedGroups    []string      `json:"feedGroups,omitempty"`
	SummaryConfig SummaryConfig `json:"summaryConfig,omitempty"`
}

type SummaryConfig struct {
	TitlesCount int `json:"titlesCount,omitempty"`
}

type HotNewsStatus struct {
	ArticlesCount  int      `json:"articlesCount,omitempty"`
	NewsLink       string   `json:"newsLink,omitempty"`
	ArticlesTitles []string `json:"articlesTitles,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
