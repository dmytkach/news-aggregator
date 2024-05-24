## Overview
The news aggregator app allows users to collect, filter and display news from various sources.
Interaction with the application occurs through the CLI.

## Detailed Description of CLI interface
1. --help.

Show all available arguments and their descriptions.

**Usage**: `go cli/cmd/main.go --help`

2. --sources

Select the desired news sources to get the news from. 

**Usage**: `go cli/cmd/main.go --sources=BBC,NBC`

3. --keywords
Specify the keywords to filter the news by.

**Usage**: `go cli/cmd/main.go --sources=BBC,NBC --keywords=Ukraine,China`

4. --date-start (--date-end)
Specify the date range to filter the news by according 
to some predefined well-known format of your choice.

**Usage**: `go cli/cmd/main.go --date-start=2024-18-05 --date-end=2024-19-05`
## Output Format
The application displays the filtered news items in the following format:
Title: <news_title>

Description: <news_description>

Link: <news_link>

Date: <publication_date>