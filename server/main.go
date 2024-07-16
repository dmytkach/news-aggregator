package main

import (
	"log"
	"net/http"
	"news-aggregator/server/handlers"
)

const (
	PORT       = ":8443"
	ServerCert = "certificates/cert.pem"
	ServerKey  = "certificates/key.pem"
)

func main() {

	http.HandleFunc("/news", handlers.News)
	http.HandleFunc("/sources", handlers.Sources)

	go handlers.FetchJob()

	log.Println("Starting server on", PORT)
	log.Fatal(http.ListenAndServeTLS(PORT, ServerCert, ServerKey, nil))
}
