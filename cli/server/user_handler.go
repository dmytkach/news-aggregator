package main

import (
	"encoding/json"
	"net/http"
	"news-aggregator/internal"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
	"news-aggregator/storage"
)

type UserHandler struct {
	s storage.Storage
}

func (h UserHandler) newsHandler(w http.ResponseWriter, r *http.Request) {
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

	v := validator.Validator{
		Sources:          sources,
		AvailableSources: h.s.AvailableSources(),
		DateStart:        dateStart,
		DateEnd:          dateEnd,
	}
	if !v.Validate() {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	a := internal.NewAggregator(
		h.s.GetAll(),
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
