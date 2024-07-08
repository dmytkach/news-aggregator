package handlers

import (
	"log"
	"news-aggregator/server/service"
	"os"
	"time"
)

var fetchInterval = time.Hour

func init() {
	intervalStr := os.Getenv("FETCH_INTERVAL")
	if intervalStr != "" {
		interval, err := time.ParseDuration(intervalStr)
		if err == nil && interval > 0 {
			fetchInterval = interval
			log.Println("Fetch interval set to ", fetchInterval)
		} else {
			log.Println("Failed to parse FETCH_INTERVAL:", intervalStr, err)
		}
	} else {
		log.Println("FETCH_INTERVAL not set")
	}
}

// FetchJob for news updating based on the set interval.
func FetchJob() {
	log.Println("Starting fetch job with interval :", fetchInterval)
	go func() {
		for {
			err := service.FetchNews()
			if err != nil {
				log.Printf("Error Fetching News: %v", err)
			}
			time.Sleep(fetchInterval)
		}
	}()
}
