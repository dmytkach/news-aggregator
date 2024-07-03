package parser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"news-aggregator/internal/entity"
	"os"
	"regexp"
	"strings"
	"time"
)

const OutputLayout = "2006-01-02"

// titleSelector to extract News titles in Usa today.
var titleSelector = "main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a"

// descriptionSelector to extract News description in Usa today.
var descriptionSelector = "data-c-br"

// dateSelector to extract News date in Usa today.
var dateSelector = "div.gnt_m_flm_sbt"

// UsaToday - parser for HTML files from Usa Today news resource.
type UsaToday struct {
	FilePath entity.PathToFile
}

func (usaTodayParser *UsaToday) CanParseFileType(ext string) bool {
	return ext == ".html"
}

// Parse - implementation of a parser for files in HTML format.
func (usaTodayParser *UsaToday) Parse() (entity.Feed, error) {
	file, err := os.Open(string(usaTodayParser.FilePath))
	if err != nil {
		return entity.Feed{}, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
			return
		}
	}(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
		return entity.Feed{}, err
	}

	baseURL := "https://www.usatoday.com"

	var allNews []entity.News
	doc.Find(titleSelector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description, _ := s.Attr(descriptionSelector)
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}
		if len(strings.TrimSpace(title)) == 0 {
			return
		}
		dateStr, _ := s.Find(dateSelector).Attr("data-c-dt")
		var parsedDate time.Time
		var err error

		if dateStr != "" {
			re := regexp.MustCompile(`[A-Za-z]+\s+\d{1,2}`)
			datePart := re.FindString(dateStr)
			if datePart != "" {
				datePart = fmt.Sprintf("%s %d", datePart, time.Now().Year())
				parsedDate, err = time.Parse("January 2 2006", datePart)
				if err != nil {
					return
				}
			}
		}

		formattedDateStr := parsedDate.Format(OutputLayout)
		formattedDate, err := time.Parse(OutputLayout, formattedDateStr)
		if formattedDate.Year() < 1000 {
			formattedDateStr = time.Now().Format(OutputLayout)
			formattedDate, err = time.Parse(OutputLayout, formattedDateStr)
		}
		allNews = append(allNews, entity.News{
			Title:       entity.Title(strings.TrimSpace(title)),
			Description: entity.Description(strings.TrimSpace(description)),
			Link:        entity.Link(strings.TrimSpace(link)),
			Date:        formattedDate,
			Source:      "usa_today",
		})

	})
	if len(allNews) == 0 {
		return entity.Feed{}, errors.New("no news found")
	}
	return entity.Feed{
		Name: "usa_today",
		News: allNews,
	}, nil
}
