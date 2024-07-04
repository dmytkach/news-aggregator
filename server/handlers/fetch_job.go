package handlers

import (
	"log"
	"news-aggregator/server/service"
	"os"
	"strconv"
	"time"
)

var fetchInterval = time.Hour

func init() {
	intervalStr := os.Getenv("FETCH_INTERVAL")
	if intervalStr != "" {
		interval, err := strconv.Atoi(intervalStr)
		if err == nil && interval > 0 {
			fetchInterval = time.Duration(interval) * time.Second
			log.Println("Fetch interval set to ", fetchInterval)
		}
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
