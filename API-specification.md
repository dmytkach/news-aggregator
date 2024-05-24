- Project Name: News Aggregator.
- Engineer name: Dmytro Tkach.
# Summary

The News Aggregator CLI is a command-line application designed to aggregate news articles from multiple sources
and provide filtering capabilities based on keywords and date ranges.

# Motivation

The purpose of this project is to create a tool that allows users to
easily access news articles from different sources and filter them based on their interests
and preferences.

Supported **use-cases**:

1) Filtering news by one source or using a set of sources.
2) Filtering news by keywords.
3) Filtering news by date range.

The News Aggregator command line interface is designed to simplify the news consumption process
by providing a centralized platform for accessing and managing news content.

# APIs design

* The `entity.News` structure defines the attributes of a news article,
  including its title, description, link, and publication date.
  This structure serves as a standardized format for representing news articles within the application.
* The `parser` package offers parsers for different file formats,
  including RSS, JSON, and HTML, allowing the application to extract news articles from various sources.
  Each parser implementation converts the raw data into structured `entity.News` objects
  for further processing.
* The `filter.NewsFilter` package provides functionality to filter news articles based on user-defined criteria,
  such as keywords and date ranges.
  This component allows users to refine their news feed to include only relevant content matching their preferences.

## Input Arguments and Output

The News Aggregator CLI accepts command-line arguments to customize the news aggregation
and filtering process. Users can specify options such as news sources, keywords, start and end dates to tailor the
results to their preferences. The output is a formatted list of news articles that match the specified criteria,
presented in a human-readable format.

## Examples

`go run cli/cmd/main.go --sources=BBC,NBC --keywords=president --date-start=2024-05-17 --date-end=2024-05-19`

This query will return news articles from the BBC and NBC, filtered by the keyword "president"
and published between May 17, 2024, and May 19, 2024, to the command line.

# Unresolved questions

1. What is the best way to work with HTML, given the unnormalized data?
2. What methods are there to analyze the file format?
3. Should the file be read only when accessed or should the files be parsed before user interaction?
