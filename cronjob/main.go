package main

import (
	"flag"
	"log"
	"news-aggregator/server/managers"
	"news-aggregator/server/service"
)

func main() {
	pathToSourcesFile := flag.String("path-to-source", "../sources.json", "Path to the file containing news sources. Default is 'server/sources.json'.")
	pathToNews := flag.String("news-folder", "../server-news/", "Path to the folder where news files are stored. Default is 'server-news/'.")

	flag.Parse()

	sourceFolder := managers.CreateSourceFolder(*pathToSourcesFile)
	newsFolder := managers.CreateNewsFolder(*pathToNews)
	urlFeed := managers.UrlFeed{}
	fetcher := service.Fetch{
		SourceManager: sourceFolder,
		NewsManager:   newsFolder,
		FeedManager:   urlFeed,
	}

	err := fetcher.UpdateNews()
	if err != nil {
		log.Printf("Error fetching news: %v", err)
	}

}
