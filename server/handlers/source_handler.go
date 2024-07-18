package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

type SourceHandler struct {
	SourceRepo managers.SourceManager
}

// Sources handles HTTP requests for managing news sources and feeds.
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// downloadSource handles HTTP POST requests to add new news feed URL.
func (sourceHandler SourceHandler) downloadSource(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	if urlStr == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	feed, err := managers.FetchFeed(entity.PathToFile(urlStr))
	if err != nil {
		log.Print("error loading feed")
		return
	}
	source, err := sourceHandler.SourceRepo.CreateSource(string(feed.Name), urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//err = sourceHandler.FetchService.UpdateNews()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(source)
}

// getSources handles HTTP GET requests to retrieve news sources.
func (sourceHandler SourceHandler) getSources(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")

	var feeds interface{}
	var err error

	if sourceName == "" {
		feeds, err = sourceHandler.SourceRepo.GetSources()
	} else {
		feeds, err = sourceHandler.SourceRepo.GetSource(sourceName)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

// updateSource handles HTTP PUT requests to update an existing news source URL.
func (sourceHandler SourceHandler) updateSource(w http.ResponseWriter, r *http.Request) {
	oldUrl := r.URL.Query().Get("oldUrl")
	newUrl := r.URL.Query().Get("newUrl")
	if oldUrl == "" || newUrl == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	err := sourceHandler.SourceRepo.UpdateSource(oldUrl, newUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// removeSource handles HTTP DELETE requests to remove a news source.
func (sourceHandler SourceHandler) removeSource(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	if sourceName == "" {
		http.Error(w, "source name is missing", http.StatusBadRequest)
		return
	}
	err := sourceHandler.SourceRepo.RemoveSourceByName(sourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
