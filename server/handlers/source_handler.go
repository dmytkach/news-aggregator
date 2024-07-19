package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/server/managers"
)

type SourceHandler struct {
	SourceManager managers.SourceManager
	FeedManager   managers.FeedManager
}

// Sources handles requests for managing news sources and feeds.
func (s SourceHandler) Sources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSources(w, r)
	case http.MethodPost:
		s.downloadSource(w, r)
	case http.MethodPut:
		s.updateSource(w, r)
	case http.MethodDelete:
		s.removeSource(w, r)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getSources handles GET requests to retrieve news sources.
func (s SourceHandler) getSources(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	log.Printf("GET request received for source: %s", sourceName)

	var feeds interface{}
	var err error

	if sourceName == "" {
		feeds, err = s.SourceManager.GetSources()
	} else {
		feeds, err = s.SourceManager.GetSource(sourceName)
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
func (s SourceHandler) downloadSource(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	log.Printf("POST request received to add source with URL: %s", urlStr)

	if urlStr == "" {
		log.Print("URL parameter is missing")
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	feed, err := s.FeedManager.Fetch(urlStr)
	if err != nil {
		log.Printf("Error loading feed from URL %s: %v", urlStr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	source, err := s.SourceManager.CreateSource(string(feed.Name), urlStr)
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
func (s SourceHandler) updateSource(w http.ResponseWriter, r *http.Request) {
	oldUrl := r.URL.Query().Get("oldUrl")
	newUrl := r.URL.Query().Get("newUrl")
	log.Printf("PUT request received to update source from URL %s to %s", oldUrl, newUrl)

	if oldUrl == "" || newUrl == "" {
		log.Print("URL parameters are missing")
		http.Error(w, "URL parameters are missing", http.StatusBadRequest)
		return
	}
	err := s.SourceManager.UpdateSource(oldUrl, newUrl)
	if err != nil {
		log.Printf("Error updating source from URL %s to %s: %v", oldUrl, newUrl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// removeSource handles DELETE requests to remove a news source.
func (s SourceHandler) removeSource(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	log.Printf("DELETE request received to remove source with name: %s", sourceName)

	if sourceName == "" {
		log.Print("Source name is missing")
		http.Error(w, "Source name is missing", http.StatusBadRequest)
		return
	}
	err := s.SourceManager.RemoveSourceByName(sourceName)
	if err != nil {
		log.Printf("Error removing source with name %s: %v", sourceName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
