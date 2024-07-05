package handlers

import (
	"encoding/json"
	"net/http"
	"news-aggregator/internal"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
)

var (
	SourceInitializer = initializers.LoadSources
)

// News handlers for HTTP GET requests to retrieve aggregated news based
// on specified query parameters.
func News(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sources := r.URL.Query().Get("sources")
	keywords := r.URL.Query().Get("keywords")
	dateStart := r.URL.Query().Get("date-start")
	dateEnd := r.URL.Query().Get("date-end")
	sortOrder := r.URL.Query().Get("sort-order")
	sortBy := r.URL.Query().Get("sort-by")
	resources, err := SourceInitializer("server-news/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	availableSources := make([]string, 0)
	for sourceName := range resources {
		availableSources = append(availableSources, sourceName)
	}
	v := validator.Validator{
		Sources:          sources,
		AvailableSources: availableSources,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
	}
	if !v.Validate() {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	a := internal.NewAggregator(
		resources,
		sources,
		initializers.InitializeFilters(&keywords, &dateStart, &dateEnd),
		sort.Options{
			Criterion: sortBy,
			Order:     sortOrder,
		})
	news, err := a.Aggregate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
