package main

import (
	"flag"
	"log"
	"net/http"
	"news-aggregator/server/handlers"
	"news-aggregator/server/managers"
	"news-aggregator/server/service"
	"time"
)

func main() {
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	port := flag.String("port", ":8443", "Port to listen on")
	serverCert := flag.String("cert", "server/certificates/cert.pem", "Path to server certificate file")
	serverKey := flag.String("key", "server/certificates/key.pem", "Path to server key file")
	fetchInterval := flag.String("fetch_interval", "1h", "Provide your fetch interval or interval will be set on 1h")
	pathToSourcesFile := flag.String("path_to_source", "server-resources/sources.json", "Please provide your source file path")
	pathToNews := flag.String("news_folder", "server-news/", "Please provide your source folder")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	sourceFolder := managers.CreateSourceFolderManager(*pathToSourcesFile)
	newsFolder := managers.CreateNewsFolderManager(*pathToNews)
	urlFeed := managers.UrlFeed{}
	sourceHandler := handlers.SourceHandler{SourceRepo: sourceFolder, FeedRepo: urlFeed}
	newsHandler := handlers.NewsHandler{NewsManager: newsFolder, SourceManager: sourceFolder}

	http.HandleFunc("/news", newsHandler.News)
	http.HandleFunc("/sources", sourceHandler.Sources)
	interval, err := time.ParseDuration(*fetchInterval)
	if err == nil && interval > 0 {
		log.Println("FeedManager interval set to ", *fetchInterval)
	} else {
		log.Println("Failed to parse FETCH_INTERVAL:", *fetchInterval, err)
	}
	job := handlers.FetchJob{Service: service.FetchService{
		SourceRepo: sourceFolder,
		NewsRepo:   newsFolder,
		Fetch:      urlFeed},
		Interval: interval}

	go job.Fetch()

	log.Println("Starting server on", *port)
	log.Fatal(http.ListenAndServeTLS(*port, *serverCert, *serverKey, nil))
}
