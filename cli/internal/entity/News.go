package entity

import (
	"os"
	"text/template"
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
func ToString(news []News) {
	var tmplFile = "cli/internal/entity/news.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, news)
	if err != nil {
		panic(err)
	}
}
