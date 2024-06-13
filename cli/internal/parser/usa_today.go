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

// selector to extract News titles in Usa today.
var selector = "main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a"

// UsaToday - parser for HTML files from Usa Today news resource.
type UsaToday struct {
	FilePath entity.PathToFile
}

// Parse - implementation of a parser for files in HTML format.
func (usaTodayParser *UsaToday) Parse() ([]entity.News, error) {
	file, err := os.Open(string(usaTodayParser.FilePath))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	baseURL := "https://www.usatoday.com"

	var allNews []entity.News
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}
		if len(strings.TrimSpace(title)) == 0 {
			return
		}
		dateStr, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
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
		})

	})
	if len(allNews) == 0 {
		return nil, errors.New("no news found")
	}
	return allNews, nil
}
