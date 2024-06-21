package main

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/internal"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
)

var resources map[string]string

func main() {
	resources = initializeDefaultResource()

	http.HandleFunc("/news", handleNews)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleNews(w http.ResponseWriter, r *http.Request) {
	sources := r.URL.Query().Get("sources")
	keywords := r.URL.Query().Get("keywords")
	sortBy := r.URL.Query().Get("sort-by")
	sortOrder := r.URL.Query().Get("sort-order")
	dateStart := r.URL.Query().Get("date-start")
	dateEnd := r.URL.Query().Get("date-end")
	v := validator.Validator{
		Sources:          sources,
		AvailableSources: resources,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
	}
	if !v.Validate() {
		return
	}
	newsFilters := internal.InitializeFilters(&keywords, &dateStart, &dateEnd)

	agg := internal.NewAggregator(resources, sources, newsFilters, sort.Options{
		Criterion: sortBy,
		Order:     sortOrder,
	})

	news := agg.Aggregate()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

// TODO: think about storage
func initializeDefaultResource() map[string]string {
	return map[string]string{
		"bbc":        "../resources/bbc-world-category-19-05-24.xml",
		"nbc":        "../resources/nbc-news.json",
		"abc":        "../resources/abcnews-international-category-19-05-24.xml",
		"washington": "../resources/washingtontimes-world-category-19-05-24.xml",
		"usatoday":   "../resources/usatoday-world-news.html",
	}
}
