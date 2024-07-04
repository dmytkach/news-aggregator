package main

import (
	"log"
	"net/http"
	"news-aggregator/server/handlers"
)

const (
	PORT        = ":8080"
	SERVER_CERT = "certificates/cert.pem"
	SERVER_KEY  = "certificates/key.pem"
)

func main() {

	http.HandleFunc("/news", handlers.News)
	http.HandleFunc("/sources", handlers.Sources)

	go handlers.FetchJob()

	log.Println("Starting server on", PORT)
	log.Fatal(http.ListenAndServeTLS(PORT, SERVER_CERT, SERVER_KEY, nil))
}
