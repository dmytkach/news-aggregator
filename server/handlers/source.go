package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/server/managers"
	"regexp"
)

type SourceHandler struct {
	SourceManager managers.SourceManager
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
	name := r.URL.Query().Get("name")
	log.Printf("POST request received to add source with Name%s ; URL: %s", name, urlStr)
	if urlStr == "" {
		log.Print("URL parameter is missing")
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	if name == "" {
		log.Print("Name parameter is missing")
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}
	reg := regexp.MustCompile(`[^\p{L}\p{N}_]+`)
	cleaned := reg.ReplaceAllString(name, "_")

	source, err := s.SourceManager.CreateSource(cleaned, urlStr)
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
	log.Printf("Successfully created source with Name: %s and URL: %s", name, urlStr)
}

// updateSource handles PUT requests to update an existing news source URL.
func (s SourceHandler) updateSource(w http.ResponseWriter, r *http.Request) {
	newUrl := r.URL.Query().Get("newUrl")
	name := r.URL.Query().Get("name")
	log.Printf("PUT request received to update source with Name %s ; New url  %s", name, newUrl)

	if newUrl == "" {
		log.Print("URL parameters are missing")
		http.Error(w, "URL parameters are missing", http.StatusBadRequest)
		return
	}
	if name == "" {
		log.Print("Name parameter is missing")
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}
	err := s.SourceManager.UpdateSource(name, newUrl)
	if err != nil {
		log.Printf("Error updating source from URL %s to %s: %v", name, newUrl, err)
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
