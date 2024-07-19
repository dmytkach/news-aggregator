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
	port := flag.String("port", ":8443", "Specify the port on which the server should listen. Default is :8443.")
	serverCert := flag.String("cert", "server/certificates/cert.pem", "Path to the server's certificate file. Default is 'server/certificates/cert.pem'.")
	serverKey := flag.String("key", "server/certificates/key.pem", "Path to the server's key file. Default is 'server/certificates/key.pem'.")
	fetchInterval := flag.String("fetch_interval", "1h", "Set the interval for fetching news updates. The default is 1 hour.")
	pathToSourcesFile := flag.String("path_to_source", "server/sources.json", "Path to the file containing news sources. Default is 'server/sources.json'.")
	pathToNews := flag.String("news_folder", "server-news/", "Path to the folder where news files are stored. Default is 'server-news/'.")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	interval, err := time.ParseDuration(*fetchInterval)
	if err == nil && interval > 0 {
		log.Println("FeedManager interval set to", *fetchInterval)
	} else {
		log.Println("Failed to parse FETCH_INTERVAL:", *fetchInterval, err)
		return
	}

	sourceFolder := managers.CreateSourceFolderManager(*pathToSourcesFile)
	newsFolder := managers.CreateNewsFolderManager(*pathToNews)
	urlFeed := managers.UrlFeed{}
	sourceHandler := handlers.SourceHandler{SourceRepo: sourceFolder, FeedRepo: urlFeed}
	newsHandler := handlers.NewsHandler{NewsManager: newsFolder, SourceManager: sourceFolder}

	http.HandleFunc("/news", newsHandler.News)
	http.HandleFunc("/sources", sourceHandler.Sources)

	job := handlers.FetchJob{Service: service.FetchService{
		SourceRepo: sourceFolder,
		NewsRepo:   newsFolder,
		Fetch:      urlFeed},
		Interval: interval}

	go job.Fetch()

	log.Println("Starting server on", *port)
	log.Fatal(http.ListenAndServeTLS(*port, *serverCert, *serverKey, nil))
}
