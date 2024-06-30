package handlers

import (
	"encoding/json"
	"net/http"
	"news-aggregator/internal"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
	"news-aggregator/server/handlers/admin/managers"
)

func NewsHandler(w http.ResponseWriter, r *http.Request) {
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
	availableSources, err := managers.GetAllSourcesNames()
	if err != nil {
		return
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
	resources, err := initializers.LoadStaticResourcesFromFolder("server-news/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	a := internal.NewAggregator(
		resources,
		sources,
		initializers.InitializeFilters(&keywords, &dateStart, &dateEnd),
		sort.Options{
			Criterion: sortBy,
			Order:     sortOrder,
		})
	news := a.Aggregate()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
