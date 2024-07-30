package entity

import (
	"time"
)

// Description of a news article.
type Description string

func (d Description) String() string {
	return string(d)
}

// Link represents URL of a news article.
type Link string

// Title of a news article.
type Title string

func (t Title) String() string {
	return string(t)
}

// News article structure with title, description, link, and date.
type News struct {
	Title       Title
	Description Description
	Link        Link
	Date        time.Time
	Source      string
}
