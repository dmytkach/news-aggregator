package main

import (
	"log"
	"net/http"
	"news-aggregator/server/handlers"
)

const PORT = ":8080"

func main() {

	http.HandleFunc("/news", handlers.News)
	http.HandleFunc("/sources", handlers.Sources)

	go handlers.FetchJob()

	log.Println("Starting server on", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
