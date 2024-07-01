package handlers

import (
	"log"
	"net/http"
	"news-aggregator/server/service"
	"time"
)

// SetInterval handles HTTP POST requests to set the fetch interval for
// automatic news updates.
func SetInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	intervalStr := r.FormValue("interval")
	if intervalStr == "" {
		http.Error(w, "Interval is required", http.StatusBadRequest)
		return
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		http.Error(w, "Invalid interval format", http.StatusBadRequest)
		return
	}

	service.SetFetchInterval(interval)

	_, err = w.Write([]byte("Interval updated successfully"))
	if err != nil {
		return
	}
}

// StartFetchScheduler news based on the set interval.
func StartFetchScheduler() {
	go func() {
		for {
			err := service.FetchNews()
			if err != nil {
				log.Fatalf("Error Fetching News: %v", err)
			}
			time.Sleep(service.GetFetchInterval())
		}
	}()
}
