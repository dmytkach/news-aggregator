package handlers

import (
	"encoding/json"
	"net/http"
	"news-aggregator/server/handlers/admin"
)

func SourcesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSourcesHandler(w, r)
	case http.MethodPost:
		downloadFeed(w, r)
	case http.MethodPut:
		updateFeed(w, r)
	case http.MethodDelete:
		removeFeed(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func downloadFeed(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	if urlStr == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	_, err := admin.NewsFeed{}.Add(urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = admin.FetchNews()
}
func getSourcesHandler(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")

	var feeds interface{}
	var err error

	if sourceName == "" {
		feeds, err = admin.NewsFeed{}.GetAll()
	} else {
		feeds, err = admin.NewsFeed{}.Get(sourceName)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)

}
func updateFeed(w http.ResponseWriter, r *http.Request) {
	oldUrl := r.URL.Query().Get("oldUrl")
	newUrl := r.URL.Query().Get("newUrl")
	if oldUrl == "" || newUrl == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	err := admin.NewsFeed{}.Update(oldUrl, newUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func removeFeed(w http.ResponseWriter, r *http.Request) {
	sourceName := r.URL.Query().Get("name")
	if sourceName == "" {
		http.Error(w, "source name is missing", http.StatusBadRequest)
		return
	}
	err := admin.NewsFeed{}.Remove(sourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
