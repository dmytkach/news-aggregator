package entity

// Feed contains multiple news articles retrieved from a specific source.
type Feed struct {
	Name SourceName
	News []News
}
