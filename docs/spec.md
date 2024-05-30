- Project Name: News Aggregator.
- Engineer name: Dmytro Tkach.

# Summary

The News Aggregator API is designed to combine news articles
from multiple sources.
The project uses parsers to extract information
from news files in different formats ( Json, Rss, Html).
In the processing process, all information is strictly structured and based
on it is possible to filter by keywords and date ranges.

# Motivation

The news aggregator API meets growing demands for effective use and news analysis.
It supports news monitoring by combining articles from different sources and offering flexible filtering options.
News Aggregator provides a convenient and effective tool that lets users know about current events.

# APIs design

## Parser:

Parsers are used in the News Aggregator API to extract data from news site sources.
Sites have unique source formats (JSON, RSS, HTML), so they are given their own parser,
responsible for the source of a certain format.
Each parser implementation converts the data into a set of strictly structured news.
P.S.  Since the content of Html files in all resources is different, 
it was decided to implement unique parsers for each resource using html
### Supported Parsers

#### 1. JSON Parser

**Description**: Parser for _JSON_ data for news extraction.

**Args**:`jsonParser.FilePath entity.PathToFile`: The path to the _JSON_ file to parse.

**Returns**:

* `[]entity.News`: A list of parsed news;
* `error`: Error object in case of failure.

**Errors**:

* `os.Open error`: Error occurred while opening the _JSON_ file.
* `json.NewDecoder error`: Error occurred while decoding _JSON_ data.

**Usage**:

```
jsonParser := parser.JsonParser{<filepath>}
news, err := jsonParser.Parse()
```

#### 2. RSS Parser

**Description**: Parses _RSS_ data to extract news articles.

**Args**:`rssParser.FilePath entity.PathToFile`: The path to the _RSS_ file to parse.

**Returns**:

* `[]entity.News`: A list of parsed news;
* `error`: Error object in case of failure.

**Errors**:

* `os.Open error`: Error occurred while opening the _RSS_ file.
* `gofeed.NewParser().Parse error`: Error occurred while parsing _RSS_ data.

**Usage**:

```
rssParser := parser.RssParser{<filepath>}
news, err := rssParser.Parse()
```


#### 3. UsaToday Parser

**Description**: Parser for _HTML_ files from Usa Today news resource.

**Args**:`rssParser.FilePath entity.PathToFile`: The path to the _HTML_ file to parse

**Returns**:

* `[]entity.News`: A list of parsed news;
* `error`: Error object in case of failure.

**Errors**:

* `os.Open error`: Error occurred while opening the _HTML_ file.
* `goquery.NewDocumentFromReader error`: Error occurred while creating a new document from the reader.

**Usage**:

```
usaTodayParser := parser.UsaTodayParser{<filepath>}
news, err := usaTodayParser.Parse()
```
## Factory method for parsers:

The factory method is used to create parser objects depending on the file provided.
It analyzes the format of the data source and selects the appropriate implementation of the parser.
### Method
**Name**: New()

**Description**: Dynamically selects and instantiates a parser object based on the provided 
data format and source.

**Args**: `entity.PathToFile`: Path to the file containing news data.

**Returns**:
`Parser`: Parser object capable of parsing the specified format.
If the specified source is not recognized returns nil.

**Usage**:

```
parser := New(source.PathToFile)
```

## News Filter
News filters are used in the news aggregator API to select news based on certain criteria.
Each filter implementation targets specific parameters to refine the content of news.
### Supported Filters

#### 1. Keyword Filter

**Description**: Filters news by specified keywords.

**Args**:`Keywords []string`: A list of keywords to filter news by.

**Returns**:`[]entity.News`: A list of news items filtered by the specified keywords.

**Usage**:

```
news []entity.news
keywordFilter := filter.KeywordFilter{Keywords: []string{"keyword1", "keyword2"}}
filteredNews := keywordFilter.Filter(news)
```
#### 2. Date Start Filter

**Description**: Filters news starting from the specified date.

**Args**:`StartDate time.Time`: The start date from which news should be filtered.

**Returns**:`[]entity.News`: A list of news items starting from the specified date.

**Usage**:

```
news []entity.news
dateStartFilter := filter.DateStartFilter{StartDate: startDate}
filteredNews := dateStartFilter.Filter(news)
```
#### 3. Date Start Filter

**Description**: Filters news up to the specified date.

**Args**:`EndDate time.Time`: The end date up to which news should be filtered.

**Returns**:`[]entity.News`: A list of news items up to the specified date.

**Usage**:

```
news []entity.news
dateEndFilter := filter.DateEndFilter{EndDate: endDate}
filteredNews := dateEndFilter.Filter(news)
```
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