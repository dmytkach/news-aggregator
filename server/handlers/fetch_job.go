package handlers

import (
	"log"
	"news-aggregator/server/service"
	"time"
)

type FetchJob struct {
	Service  service.FetchService
	Interval time.Duration
}

// Fetch for news updating based on the set interval.
func (f FetchJob) Fetch() {
	log.Println("Starting fetch job with interval :", f.Interval)
	go func(fetchInterval time.Duration) {
		for {
			err := f.Service.UpdateNews()
			if err != nil {
				log.Printf("Error Fetching News: %v", err)
			}
			time.Sleep(fetchInterval)
		}
	}(f.Interval)
}
