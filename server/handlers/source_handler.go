package handlers

import (
	"encoding/json"
	"net/http"
	"news-aggregator/server/service"
)

// Sources handles HTTP requests for managing news sources and feeds.
func Sources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSources(w, r)
	case http.MethodPost:
		downloadSource(w, r)
	case http.MethodPut:
		updateSource(w, r)
	case http.MethodDelete:
		removeSource(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// downloadSource handles HTTP POST requests to add new news feed URL.
func downloadSource(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	if urlStr == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	source, err := service.AddSource(urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = service.FetchNews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(source)
}

// getSources handles HTTP GET requests to retrieve news sources.
func getSources(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")

	var feeds interface{}
	var err error

	if sourceName == "" {
		feeds, err = service.GetSources()
	} else {
		feeds, err = service.GetSource(sourceName)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)

}

// updateSource handles HTTP PUT requests to update an existing news source URL.
func updateSource(w http.ResponseWriter, r *http.Request) {
	oldUrl := r.URL.Query().Get("oldUrl")
	newUrl := r.URL.Query().Get("newUrl")
	if oldUrl == "" || newUrl == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	err := service.UpdateSource(oldUrl, newUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// removeSource handles HTTP DELETE requests to remove a news source.
func removeSource(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	if sourceName == "" {
		http.Error(w, "source name is missing", http.StatusBadRequest)
		return
	}
	err := service.RemoveSource(sourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
