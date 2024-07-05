package managers

import (
	"encoding/json"
	"fmt"
	"news-aggregator/internal/entity"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestAddNews(t *testing.T) {
	mockTime := time.Date(2024, 7, 3, 12, 0, 0, 0, time.UTC)
	news := []entity.News{
		{
			Title:       "Test News 1",
			Description: "This is a test news article 1",
			Link:        "https://example.com/news1",
			Date:        mockTime,
			Source:      "test_source",
		},
		{
			Title:       "Test News 2",
			Description: "This is a test news article 2",
			Link:        "https://example.com/news2",
			Date:        mockTime,
			Source:      "test_source",
		},
	}
	NewsFolder = "test-data"
	newsSource := "test-source"
	err := AddNews(news, newsSource)
	if err != nil {
		t.Fatalf("AddNews() failed: %v", err)
	}

	finalFileName := fmt.Sprintf("%s/%s.json", newsSource, timeNow)
	finalFilePath := filepath.Join(NewsFolder, finalFileName)

	fileData, err := os.ReadFile(finalFilePath)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}

	var storedNews []entity.News
	err = json.Unmarshal(fileData, &storedNews)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON data: %v", err)
	}

	if !reflect.DeepEqual(storedNews, news) {
		t.Errorf("Stored news data does not match expected data. Got: %v, Expected: %v", storedNews, news)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(NewsFolder)
		if err != nil {
			t.Errorf("Failed to clean up server-news folder: %v", err)
		}
	})
}

func TestGetNewsFromFile(t *testing.T) {
	existingFilePath := "../../internal/testdata/news.json"

	got, err := GetNewsFromFile(existingFilePath)
	if err != nil {
		t.Fatalf("GetNewsFromFile() failed: %v", err)
	}

	wants := entity.Feed{
		Name: "nbc_news",
		News: []entity.News{
			{
				Title:       "Container ship that struck Baltimore bridge will be removed from the site 'within days,' Maryland governor says",
				Description: "Ship that struck Francis Scott Key Bridge in Baltimore will be removed \"within days,\" Maryland Gov. Wes Moore says",
				Link:        "https://www.nbcnews.com/politics/politics-news/francis-scott-key-bridge-ship-removal-wes-moore-baltimore-rcna152955",
				Date:        time.Date(2024, 5, 19, 14, 6, 47, 0, time.UTC),
				Source:      "nbc_news",
			},
			{
				Title:       "Harris says more Indian American representation is needed in government",
				Description: "Addressing a crowd of Indian Americans this week, Vice President Kamala Harris asserted the importance of voting and running. But Biden and Harris approval among the group has fallen.",
				Link:        "https://www.nbcnews.com/news/asian-america/kamala-harris-more-indian-american-representation-needed-government-rcna152761",
				Date:        time.Date(2024, 5, 17, 19, 48, 19, 0, time.UTC),
				Source:      "nbc_news",
			},
			{
				Title:       "Atlanta officer accused of killing Lyft driver allegedly said victim was ‘gay fraternity’ recruiter",
				Description: "An Atlanta police officer accused of murdering a Lyft driver allegedly said the victim was in a gay fraternity trying to recruit him.",
				Link:        "https://www.nbcnews.com/nbc-out/out-news/atlanta-officer-accused-killing-lyft-driver-allegedly-said-victim-was-rcna152751",
				Date:        time.Date(2024, 5, 17, 14, 29, 43, 0, time.UTC),
				Source:      "nbc_news",
			},
		}}

	if !reflect.DeepEqual(got, wants) {
		t.Errorf("Retrieved news data does not match expected data. Got: %v, Expected: %v", got, wants)
	}
}
func TestGetNewsFromFolder(t *testing.T) {
	NewsFolder = "../../internal/testdata"
	got, err := GetNewsFromFolder("bbc_news")
	if err != nil {
		t.Fatalf("GetNewsFromFile() failed: %v", err)
	}

	wants := []entity.News{
		{Title: "Why parents are locking themselves in cells at Korean 'happiness factory'",
			Description: "Some South Koreans are spending time in a cell to try to understand their socially isolated children.",
			Link:        "https://www.bbc.com/news/articles/c2x0le06kn7o",
			Date:        time.Date(2024, time.June, 30, 1, 6, 19, 0, time.UTC),
			Source:      "bbc_news",
		},
		{
			Title:       "Watch England fans go wild as Bellingham scores late equaliser",
			Description: "Watch England fans erupt at a fanpark in Wembley as Jude Bellingham scores a late stunner to send the game to extra time against Slovakia in the Euro 2024 last-16 match in Gelsenkirchen.\n\n\n",
			Link:        "https://www.bbc.com/sport/football/videos/cl4yj1ve5z7o",
			Source:      "bbc_news",
			Date:        time.Date(2024, 6, 30, 19, 31, 26, 0, time.UTC),
		},
		{
			Title:       "Sunak and Labour's Final Sunday Pitch",
			Description: "The PM and Labour’s campaign chief are in Laura’s BBC One studio",
			Link:        "https://www.bbc.co.uk/sounds/play/p0j7g5dz",
			Date:        time.Date(2024, 6, 30, 13, 52, 0, 0, time.UTC),
			Source:      "bbc_news",
		},
	}

	if !reflect.DeepEqual(got, wants) {
		t.Errorf("Retrieved news data does not match expected data. Got: %v, Expected: %v", got, wants)
	}
}
