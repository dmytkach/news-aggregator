package handlers

import (
	"log"
	"news-aggregator/server/service"
	"time"
)

type FetchJob struct {
	FetchService  service.FetchService
	FetchInterval time.Duration
}

// Fetch for news updating based on the set interval.
func (f FetchJob) Fetch() {
	log.Println("Starting fetch job with interval :", f.FetchInterval)
	go func() {
		for {
			err := f.FetchService.UpdateNews()
			if err != nil {
				log.Printf("Error Fetching News: %v", err)
			}
			time.Sleep(f.FetchInterval)
		}
	}()
}
