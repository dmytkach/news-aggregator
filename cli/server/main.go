package main

import (
	"log"
	"net/http"
	"news-aggregator/server/handlers"
)

func main() {

	http.HandleFunc("/news", handlers.NewsHandler)
	http.HandleFunc("/sources", handlers.SourcesHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
