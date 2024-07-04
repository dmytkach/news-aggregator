## Overview
The news aggregator app allows users to collect, filter and display news from various sources.
Interaction with the application occurs through HTTP endpoints and a CLI interface.

## Server API

The server API exposes the following endpoints:

### `/news`

Allows users to retrieve aggregated news based on various query parameters.

#### Query Parameters

- `sources`: (Optional) Comma-separated list of news sources from which to fetch news.
- `keywords`: (Optional) Comma-separated list of keywords to filter news articles.
- `date-start`: (Optional) Start date to filter news articles. Should be in `YYYY-MM-DD` format.
- `date-end`: (Optional) End date to filter news articles. Should be in `YYYY-MM-DD` format.
- `sort-order`: (Optional) Specifies the order in which news articles should be sorted. Options: `asc` (ascending) or `desc` (descending).
- `sort-by`: (Optional) Specifies the criterion for sorting news articles. Options may include `date`, `title`, or other relevant criteria depending on your implementation.

#### Example Usage
```
GET /news?sources=BBC,CNN&keywords=technology,science&date-start=2024-06-01&date-end=2024-06-30&sort-order=desc&sort-by=date
```
### `/sources`
Managing news sources including adding, updating, and removing news sources.

**Supported Methods**:
  - `GET`: Retrieves information about news sources.
  - `POST`: Adds a new news source.
  - `PUT`: Updates an existing news source.
  - `DELETE`: Removes a news source.
### Starting the Server

When you start the server, you can set the news fetching interval
using the `FETCH_INTERVAL` environment variable. The value must be specified in seconds.
#### Example usage:
```
go run server/main.go
```
## CLI API
1. --help.

Show all available arguments and their descriptions.

**Usage**: `go cli/main.go --help`

2. --sources

Select the desired news sources to get the news from. 

**Usage**: `go cli/main.go --sources=BBC,NBC`

3. --keywords
Specify the keywords to filter the news by.

**Usage**: `go cli/main.go --sources=BBC,NBC --keywords=Ukraine,China`

4. --date-start (--date-end)
Specify the date range to filter the news by according 
to some predefined well-known format of your choice.

**Usage**: `go cli/main.go --date-start=2024-18-05 --date-end=2024-19-05`
## Output Format
The application displays the filtered news items in the following format:
Title: <news_title>

Description: <news_description>

Link: <news_link>

Date: <publication_date>

Source: <news_source>