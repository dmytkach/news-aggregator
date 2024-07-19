package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/internal"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
	"news-aggregator/server/managers"
)

type NewsHandler struct {
	NewsManager   managers.NewsManager
	SourceManager managers.SourceManager
}

// News handler for GET requests to retrieve aggregated news based
// on specified query parameters.
func (newsHandler NewsHandler) News(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sources := r.URL.Query().Get("sources")
	keywords := r.URL.Query().Get("keywords")
	dateStart := r.URL.Query().Get("date-start")
	dateEnd := r.URL.Query().Get("date-end")
	sortOrder := r.URL.Query().Get("sort-order")
	sortBy := r.URL.Query().Get("sort-by")

	log.Printf("Received GET request with parameters - Sources: %s, Keywords: %s, DateStart: %s, DateEnd: %s, SortOrder: %s, SortBy: %s",
		sources, keywords, dateStart, dateEnd, sortOrder, sortBy)

	s, err := newsHandler.SourceManager.GetSources()
	if err != nil {
		log.Printf("Error retrieving sources: %v", err)
		http.Error(w, "Error retrieving news source file paths", http.StatusInternalServerError)
		return
	}

	availableSources := make([]string, 0)
	for _, source := range s {
		availableSources = append(availableSources, string(source.Name))
	}
	log.Printf("Available sources: %v", availableSources)

	resources, err := newsHandler.NewsManager.GetNewsSourceFilePath(availableSources)
	if err != nil {
		log.Printf("Error getting news source file paths: %v", err)
		http.Error(w, "Error retrieving news source file paths", http.StatusBadRequest)
		return
	}

	sortOptions := sort.Options{
		Criterion: sortBy,
		Order:     sortOrder,
	}
	config := validator.Config{
		Sources:          sources,
		AvailableSources: availableSources,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
		SortOptions:      sortOptions,
	}

	v := validator.NewValidator(config)
	if !v.Validate() {
		log.Printf("Invalid query parameters: %v", config)
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	a := internal.NewAggregator(
		resources,
		sources,
		initializers.InitializeFilters(&keywords, &dateStart, &dateEnd),
		sortOptions)
	news, err := a.Aggregate()
	if err != nil {
		log.Printf("Error aggregating news: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(news)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
