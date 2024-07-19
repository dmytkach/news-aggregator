## Overview

The news aggregator app allows users to collect, filter and display news from various sources.
Interaction with the application occurs through HTTPs endpoints and a CLI interface.

## Server API

The server API exposes the following endpoints:

### `/news`

Allows users to retrieve aggregated news based on various query parameters.

#### Query Parameters

- `sources`: (Optional) Comma-separated list of news sources from which to fetch news.
- `keywords`: (Optional) Comma-separated list of keywords to filter news articles.
- `date-start`: (Optional) Start date to filter news articles. Should be in `YYYY-MM-DD` format.
- `date-end`: (Optional) End date to filter news articles. Should be in `YYYY-MM-DD` format.
- `sort-order`: (Optional) Specifies the order in which news articles should be sorted. Options: `asc` (ascending)
  or `desc` (descending).
- `sort-by`: (Optional) Specifies the criterion for sorting news articles. Options may include `date`, `title`, or other
  relevant criteria depending on your implementation.

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

When you start the server, you can configure various settings using command-line flags.
Here are the available flags and their descriptions:

1. --help:

Displays all available arguments and their descriptions.
If you run the server with this flag, it will show a help message with details about each flag.

**Usage**: `go run server/main.go --help`

2. --port:

Specifies the port on which the server will listen. The default value is :8443.
You can set it to a different port if needed.

**Usage**: `go run server/main.go --port=:443`

3. --cert:

Provides the path to the server certificate file. The default path is server/certificates/cert.pem.
If you have a different certificate file, specify its path here.

**Usage**:`go run server/main.go --cert=/path/to/your/cert.pem`

4. --key:

Provides the path to the server key file. The default path is server/certificates/key.pem.
If you have a different key file, specify its path here.

**Usage**:`go run server/main.go --key=/path/to/your/key.pem`

5. --fetch_interval:

Sets the interval for fetching news updates. The default value is 1h (1 hour).
You can specify a different interval using valid time units (e.g., 30s for 30 seconds, 5m for 5 minutes).

**Usage**: `go run server/main.go --fetch_interval=30s`

6. --path_to_source:

Specifies the path to the source file containing news sources.
The default path is server/sources.json.
If your sources file is located elsewhere, provide its path here.

**Usage**: `go run server/main.go --path_to_source=/path/to/your/sources.json`

7. --news_folder:

Specifies the folder where news files are stored. The default folder is server-news/.
If your news folder is located elsewhere, provide its path here.

**Usage**: `go run server/main.go --news_folder=/path/to/your/news_folder`

## Docker Instructions

This project provides a Docker image for the news aggregator application. Below are the instructions for using Docker
with this project.

### 1. Pull the Docker Image

To pull the Docker image for this project, use the following command. Replace `news-aggregator` and `latest` with the
appropriate image name and tag if needed.

```
docker pull news-aggregator:latest
````

### 2. Build the Docker Image

If you want to build the Docker image yourself, follow these steps:

Make sure you have Docker installed on your machine.
Navigate to the root directory of the project.
Build the Docker image using the following command:

```
docker build -t news-aggregator:latest .
```

### 3. Run the Docker Container

To run the Docker container, use the following command:

```
docker run --rm -p 8443:8443 news-aggregator:latest -fetch_interval=30s 
    -cert=/path/to/your/cert.pem \
    -key=/path/to/your/key.pem \
    -path_to_source=sources.json \
    -news_folder=server-news/
```

## Instructions for starting the server from CLI:
This project also provides a CLI interface for the news aggregator server.
To start the server via CLI, use the following command variations:

### Custom settings:

```
go run server/main.go 
    --port=:8080
    --cert=/path/to/your/cert.pem
    --key=/path/to/your/key.pem --fetch_interval=30s
    --path_to_source=/path/to/your/sources.json
    --news_folder=/path/to/your/news_folder
```

### Default settings:

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