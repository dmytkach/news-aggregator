package filters

import (
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/internal/entity"
	"slices"
	"strings"
)

// Keyword filters news by keywords.
type Keyword struct {
	Keywords []string
}

// Filter news by keywords in the title and description.
func (k *Keyword) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	keywords := getStemKeywords(k)
	for _, item := range news {
		titles := strings.Split(strings.ToLower(string(item.Title)), " ")
		description := strings.Split(strings.ToLower(string(item.Description)), " ")
		for _, stemmedKeyword := range keywords {
			if slices.Contains(titles, stemmedKeyword) || slices.Contains(description, stemmedKeyword) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func getStemKeywords(k *Keyword) []string {
	var stemmedWords = make([]string, 0)
	for _, keyword := range k.Keywords {
		stemmedWords = append(stemmedWords, strings.ToLower(keyword))          // add original keyword
		stemmedWords = append(stemmedWords, porterstemmer.StemString(keyword)) // add stemmed keyword
	}
	return stemmedWords
}
func (k *Keyword) String() string {
	return "keywords=" + strings.Join(k.Keywords, ",")
}
