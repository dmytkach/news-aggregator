- Project Name: News Aggregator.
- Engineer name: Dmytro Tkach.

# Summary

The News Aggregator CLI is a command-line application designed to aggregate news articles from multiple sources
and provide filtering capabilities based on keywords and date ranges.

### Supported **use-cases**:

1) Filtering news by one source or using a set of sources.
2) Filtering news by keywords.
3) Filtering news by date range.

# APIs design

The client level works with the user entering parameters into the program - sources, keywords, dates.
News sources are processed separately. For each source, the format of the file used is analyzed and the file
is processed by the corresponding parser.
Processed news is converted to the News type and stored in the program at runtime.
If the process of processing and saving news is successful,
the filtering process begins in accordance with the set of entered parameters.
Filters use a single interface. For each parameter type,
an implementation of this interface is created.
After which the user receives a list of all news that satisfies the request.

* The `entity` package provides APIs for the `News` and `Resource` structures.
  The `News` structure is a standardized format for presenting news articles in an application.
  The `Resource` structure displays the name of the resource and the corresponding file path for that resource.
* The `parser` package provides an API for processing different file formats,
  including RSS, JSON and HTML.
  Each parser implementation converts raw data into structured `entity.News` objects.
* The `filter` package provides a API for filtering `entity.News` based on valid criteria.
  Keywords and date ranges are supported.

## Examples

### Command line request:

`go run cli/cmd/main.go --sources=BBC --keywords=president --date-start=2024-05-17`

This query will return news articles from the BBC , filtered by the keyword "president"
and published after May 17, 2024, to the command line.

### Request processing process:

`resources = append(resources, entity.Resource{Name: "BBC", PathToFile: "../.xml"})` - define resources.

`parser.GetParser(source.PathToFile).Parse()` - get a list of processed news for this file.

`var filters []filter.NewsFilter` - create an empty list of filters.

`filters = append(filters, &filter.KeywordFilter{Keywords: keywordList})` - add KeywordFilter.

`filters = append(filters, &filter.DateStartFilter{StartDate: startDate})` - add a date filter.

`result = newsFilter.Filter(result)` - use the general interface for filtering news and save for output.

# Unresolved questions

1. What is the best way to work with HTML, given the unnormalized data?