package main

import (
	"log"
	"net/http"
	"news-aggregator/server/handlers"
)

func main() {

	http.HandleFunc("/news", handlers.News)
	http.HandleFunc("/sources", handlers.Sources)
	http.HandleFunc("/set-interval", handlers.SetInterval)

	go handlers.StartFetchScheduler()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
