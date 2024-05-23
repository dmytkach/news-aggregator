package entity

import (
	"strings"
	"time"
)

// Description of a news article.
type Description string

// Link represents URL of a news article.
type Link string

// Title of a news article.
type Title string

// News article structure with title, description, link, and date.
type News struct {
	Title       Title
	Description Description
	Link        Link
	Date        time.Time
}

// ToString formats the news article into a human-readable string.
func (news *News) ToString() string {
	parts := []string{
		string("Title: " + news.Title),
		string("Description: " + news.Description),
		string("Link: " + news.Link),
		"Date: " + news.Date.String(),
	}
	return strings.Join(parts, "\n")
}
