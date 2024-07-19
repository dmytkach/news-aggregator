package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/server/managers"
)

type SourceHandler struct {
	SourceRepo managers.SourceManager
	FeedRepo   managers.FeedManager
}

// Sources handles requests for managing news sources and feeds.
func (sourceHandler SourceHandler) Sources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sourceHandler.getSources(w, r)
	case http.MethodPost:
		sourceHandler.downloadSource(w, r)
	case http.MethodPut:
		sourceHandler.updateSource(w, r)
	case http.MethodDelete:
		sourceHandler.removeSource(w, r)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getSources handles GET requests to retrieve news sources.
func (sourceHandler SourceHandler) getSources(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	log.Printf("GET request received for source: %s", sourceName)

	var feeds interface{}
	var err error

	if sourceName == "" {
		feeds, err = sourceHandler.SourceRepo.GetSources()
	} else {
		feeds, err = sourceHandler.SourceRepo.GetSource(sourceName)
	}

	if err != nil {
		log.Printf("Error retrieving sources: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// downloadSource handles POST requests to add new news feed URL.
func (sourceHandler SourceHandler) downloadSource(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	log.Printf("POST request received to add source with URL: %s", urlStr)

	if urlStr == "" {
		log.Print("URL parameter is missing")
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	feed, err := sourceHandler.FeedRepo.Fetch(urlStr)
	if err != nil {
		log.Printf("Error loading feed from URL %s: %v", urlStr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	source, err := sourceHandler.SourceRepo.CreateSource(string(feed.Name), urlStr)
	if err != nil {
		log.Printf("Error creating source for URL %s: %v", urlStr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(source); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// updateSource handles PUT requests to update an existing news source URL.
func (sourceHandler SourceHandler) updateSource(w http.ResponseWriter, r *http.Request) {
	oldUrl := r.URL.Query().Get("oldUrl")
	newUrl := r.URL.Query().Get("newUrl")
	log.Printf("PUT request received to update source from URL %s to %s", oldUrl, newUrl)

	if oldUrl == "" || newUrl == "" {
		log.Print("URL parameters are missing")
		http.Error(w, "URL parameters are missing", http.StatusBadRequest)
		return
	}
	err := sourceHandler.SourceRepo.UpdateSource(oldUrl, newUrl)
	if err != nil {
		log.Printf("Error updating source from URL %s to %s: %v", oldUrl, newUrl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// removeSource handles DELETE requests to remove a news source.
func (sourceHandler SourceHandler) removeSource(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	log.Printf("DELETE request received to remove source with name: %s", sourceName)

	if sourceName == "" {
		log.Print("Source name is missing")
		http.Error(w, "Source name is missing", http.StatusBadRequest)
		return
	}
	err := sourceHandler.SourceRepo.RemoveSourceByName(sourceName)
	if err != nil {
		log.Printf("Error removing source with name %s: %v", sourceName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
