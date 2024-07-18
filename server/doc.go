// Package main provides API for work with the news aggregator server.
// This server aggregates news from various sources and provides endpoints
// for managing news sources and fetching aggregated news based on query parameters.
//
// Starting the Server:
// The server starts on port 8443 and exposes two main endpoints:
//   - /news: Endpoint for fetching aggregated news.
//   - /sources: Endpoint for managing news sources (add, update, delete).
//
// The server initiates a background fetch job based on the interval set by
// the FETCH_INTERVAL environment variable (in seconds). By default, it fetches news every hour.
package main
