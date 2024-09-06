package managers

import (
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"reflect"
	"testing"
	"time"
)

func TestUrlFeed_Fetch(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0">
			<channel>
				<title>Mock News</title>
				<item>
					<title>Mock Title</title>
					<description>Mock Description</description>
					<link>https://mock.link</link>
					<pubDate>Mon, 01 Jan 2001 00:00:00 +0000</pubDate>
				</item>
			</channel>
		</rss>
		`))
	}))
	defer mockServer.Close()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			url:     mockServer.URL,
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "http://invalid.url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UrlFeed{}
			_, err := f.FetchFeed(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlFeed.FeedManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestGetFeedFromFile(t *testing.T) {
	existingFilePath := "../../internal/testdata/news.json"
	invalidFilePath := "../../internal/testdata/invalid.json"

	tests := []struct {
		name     string
		filePath string
		want     []entity.News
		wantErr  bool
	}{
		{
			name:     "Valid File",
			filePath: existingFilePath,
			want: []entity.News{
				{
					Title:       "Container ship that struck Baltimore bridge will be removed from the site 'within days,' Maryland governor says",
					Description: "Ship that struck Francis Scott Key Bridge in Baltimore will be removed \"within days,\" Maryland Gov. Wes Moore says",
					Link:        "https://www.nbcnews.com/politics/politics-news/francis-scott-key-bridge-ship-removal-wes-moore-baltimore-rcna152955",
					Date:        time.Date(2024, 5, 19, 14, 6, 47, 0, time.UTC),
					Source:      "NBC News",
				},
				{
					Title:       "Harris says more Indian American representation is needed in government",
					Description: "Addressing a crowd of Indian Americans this week, Vice President Kamala Harris asserted the importance of voting and running. But Biden and Harris approval among the group has fallen.",
					Link:        "https://www.nbcnews.com/news/asian-america/kamala-harris-more-indian-american-representation-needed-government-rcna152761",
					Date:        time.Date(2024, 5, 17, 19, 48, 19, 0, time.UTC),
					Source:      "NBC News",
				},
				{
					Title:       "Atlanta officer accused of killing Lyft driver allegedly said victim was ‘gay fraternity’ recruiter",
					Description: "An Atlanta police officer accused of murdering a Lyft driver allegedly said the victim was in a gay fraternity trying to recruit him.",
					Link:        "https://www.nbcnews.com/nbc-out/out-news/atlanta-officer-accused-killing-lyft-driver-allegedly-said-victim-was-rcna152751",
					Date:        time.Date(2024, 5, 17, 14, 29, 43, 0, time.UTC),
					Source:      "NBC News",
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid File",
			filePath: invalidFilePath,
			want:     []entity.News{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFeedFromFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFeedFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFeedFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
