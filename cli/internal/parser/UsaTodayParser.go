package parser

import (
	"NewsAggregator/cli/internal/entity"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
	"time"
)

// UsaTodayParser - parser for HTML files from Usa Today news resource.
type UsaTodayParser struct {
	FilePath entity.PathToFile
}

// Parse - implementation of a parser for files in HTML format.
func (usaTodayParser *UsaTodayParser) Parse() ([]entity.News, error) {
	file, err := os.Open(string(usaTodayParser.FilePath))
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			print(err)
		}
	}(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	const layout = "January 2, 2006"
	const outputLayout = "2006-01-02"
	baseURL := "https://www.usatoday.com"

	var news []entity.News
	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		dateStr, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Println("Error parsing date:", err)
			parsedDate = time.Time{}
		}

		formattedDateStr := parsedDate.Format(outputLayout)
		formattedDate, err := time.Parse(outputLayout, formattedDateStr)
		if err != nil {
			log.Println("Error formatting date:", err)
		}

		news = append(news, entity.News{
			Title:       entity.Title(strings.TrimSpace(title)),
			Description: entity.Description(strings.TrimSpace(description)),
			Link:        entity.Link(strings.TrimSpace(link)),
			Date:        formattedDate,
		})
	})

	return news, nil
}
