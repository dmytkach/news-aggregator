package main

import (
	"flag"
	"log"
	"net/http"
	"news-aggregator/server/handlers"
	"news-aggregator/server/managers"
)

// main initializes and starts the news aggregator server.
func main() {
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	port := flag.String("port", ":8443", "Specify the port on which the server should listen. Default is :8443.")
	pathToSourcesFile := flag.String("path-to-source", "server/sources.json", "Path to the file containing news sources. Default is 'server/sources.json'.")
	pathToNews := flag.String("news-folder", "server-news/", "Path to the folder where news files are stored. Default is 'server-news/'.")
	certFile := flag.String("tls-cert", "/etc/tls/certs/tls.crt", "Path to the TLS certificate file.")
	keyFile := flag.String("tls-key", "/etc/tls/certs/tls.key", "Path to the TLS key file.")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	sourceFolder := managers.CreateSourceFolder(*pathToSourcesFile)
	newsFolder := managers.CreateNewsFolder(*pathToNews)
	sourceHandler := handlers.SourceHandler{SourceManager: sourceFolder}
	newsHandler := handlers.NewsHandler{NewsManager: newsFolder, SourceManager: sourceFolder}

	http.HandleFunc("/news", newsHandler.News)
	http.HandleFunc("/sources", sourceHandler.Sources)

	log.Println("Starting server on", *port)
	err := http.ListenAndServeTLS(*port, *certFile, *keyFile, nil)
	if err != nil {
		log.Print("Error starting server: ", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
