// package managers
//
// import (
//
//	"encoding/json"
//	"fmt"
//	"news-aggregator/internal/entity"
//	"os"
//	"path/filepath"
//	"reflect"
//	"testing"
//	"time"
//
// )
//
//	func TestAddNews(t *testing.T) {
//		NewsFolder, expectedNews := setupTestData(t)
//		defer cleanupTestData(t, NewsFolder)
//
//		finalFileName := fmt.Sprintf("test-source/%s.json", timeNow)
//		finalFilePath := filepath.Join(NewsFolder, finalFileName)
//
//		fileData, err := os.ReadFile(finalFilePath)
//		if err != nil {
//			t.Fatalf("Failed to read created file: %v", err)
//		}
//
//		var storedNews []entity.News
//		err = json.Unmarshal(fileData, &storedNews)
//		if err != nil {
//			t.Fatalf("Failed to unmarshal JSON data: %v", err)
//		}
//
//		if !reflect.DeepEqual(storedNews, expectedNews) {
//			t.Errorf("Stored news data does not match expected data. Got: %v, Expected: %v", storedNews, expectedNews)
//		}
//
//		invalidNewsHandler := newsFolder{path: string([]byte{0x00})}
//		err = invalidNewsHandler.AddNews(expectedNews, "invalid-source")
//		if err == nil {
//			t.Errorf("Expected an error when creating a directory with an invalid path, but got none")
//		}
//
//		invalidFilePath := filepath.Join("?'@%=", finalFileName)
//		err = os.MkdirAll(filepath.Dir(invalidFilePath), 0755)
//		if err == nil {
//			t.Fatalf("Expected an error when creating a directory with an invalid path, but got none")
//		}
//
//		err = os.WriteFile(invalidFilePath, []byte("test"), 0644)
//		if err == nil {
//			t.Errorf("Expected an error when writing to a file with an invalid path, but got none")
//		}
//	}
//
//	func setupTestData(t *testing.T) (string, []entity.News) {
//		t.Helper()
//
//		mockTime := time.Date(2024, 7, 3, 12, 0, 0, 0, time.UTC)
//		news := []entity.News{
//			{
//				Title:       "Test News 1",
//				Description: "This is a test news article 1",
//				Link:        "https://example.com/news1",
//				Date:        mockTime,
//				Source:      "test_source",
//			},
//			{
//				Title:       "Test News 2",
//				Description: "This is a test news article 2",
//				Link:        "https://example.com/news2",
//				Date:        mockTime,
//				Source:      "test_source",
//			},
//		}
//
//		NewsFolder := "test-data"
//		newsHandler := newsFolder{path: NewsFolder}
//		err := newsHandler.AddNews(news, "test-source")
//		if err != nil {
//			t.Fatalf("AddNews() failed: %v", err)
//		}
//
//		return NewsFolder, news
//	}
//
//	func cleanupTestData(t *testing.T, path string) {
//		t.Helper()
//
//		err := os.RemoveAll(path)
//		if err != nil {
//			t.Errorf("Failed to clean up test data folder: %v", err)
//		}
//	}
//
// //func TestAddNews(t *testing.T) {
// //	NewsFolder, expectedNews := setupTestData(t)
// //	defer cleanupTestData(t, NewsFolder)
// //
// //	finalFileName := fmt.Sprintf("test-source/%s.json", timeNow)
// //	finalFilePath := filepath.Join(NewsFolder, finalFileName)
// //
// //	fileData, err := os.ReadFile(finalFilePath)
// //	if err != nil {
// //		t.Errorf("Failed to read created file: %v", err)
// //	}
// //
// //	var storedNews []entity.News
// //	err = json.Unmarshal(fileData, &storedNews)
// //	if err != nil {
// //		t.Errorf("Failed to unmarshal JSON data: %v", err)
// //	}
// //
// //	if !reflect.DeepEqual(storedNews, expectedNews) {
// //		t.Errorf("Stored news data does not match expected data. Got: %v, Expected: %v", storedNews, expectedNews)
// //	}
// //
// //	invalidNewsHandler := newsFolder{path: string([]byte{0x00})}
// //	err = invalidNewsHandler.AddNews(expectedNews, "invalid-source")
// //	if err == nil {
// //		t.Errorf("Expected an error when creating a directory with an invalid path, but got none")
// //	}
// //
// //	invalidFilePath := filepath.Join("?'@%=", finalFileName)
// //	err = os.MkdirAll(filepath.Dir(invalidFilePath), 0755)
// //	if err == nil {
// //		t.Fatalf("Expected an error when creating a directory with an invalid path, but got none")
// //	}
// //
// //	err = os.WriteFile(invalidFilePath, []byte("test"), 0644)
// //	if err == nil {
// //		t.Errorf("Expected an error when writing to a file with an invalid path, but got none")
// //	}
// //}
//
//	func TestGetNewsFromFolder(t *testing.T) {
//		NewsFolder := "../../internal/testdata"
//		newsHandler := newsFolder{path: NewsFolder}
//		got, err := newsHandler.GetNewsFromFolder("bbc_news")
//		if err != nil {
//			t.Fatalf("GetNewsFromFolder() failed: %v", err)
//		}
//
//		wants := []entity.News{
//			{
//				Title:       "Why parents are locking themselves in cells at Korean 'happiness factory'",
//				Description: "Some South Koreans are spending time in a cell to try to understand their socially isolated children.",
//				Link:        "https://www.bbc.com/news/articles/c2x0le06kn7o",
//				Date:        time.Date(2024, time.June, 30, 1, 6, 19, 0, time.UTC),
//				Source:      "bbc_news",
//			},
//			{
//				Title:       "Watch England fans go wild as Bellingham scores late equaliser",
//				Description: "Watch England fans erupt at a fanpark in Wembley as Jude Bellingham scores a late stunner to send the game to extra time against Slovakia in the Euro 2024 last-16 match in Gelsenkirchen.\n\n\n",
//				Link:        "https://www.bbc.com/sport/football/videos/cl4yj1ve5z7o",
//				Date:        time.Date(2024, 6, 30, 19, 31, 26, 0, time.UTC),
//				Source:      "bbc_news",
//			},
//			{
//				Title:       "Sunak and Labour's Final Sunday Pitch",
//				Description: "The PM and Labour’s campaign chief are in Laura’s BBC One studio",
//				Link:        "https://www.bbc.co.uk/sounds/play/p0j7g5dz",
//				Date:        time.Date(2024, 6, 30, 13, 52, 0, 0, time.UTC),
//				Source:      "bbc_news",
//			},
//		}
//
//		if !reflect.DeepEqual(got, wants) {
//			t.Errorf("Retrieved news data does not match expected data. Got: %v, Expected: %v", got, wants)
//		}
//		errorFolder := "test-error-data"
//		err = os.MkdirAll(filepath.Join(NewsFolder, errorFolder), 0755)
//		if err != nil {
//			t.Fatalf("Failed to create test erroneous folder: %v", err)
//		}
//		_, err = os.Create(filepath.Join(NewsFolder, errorFolder, "invalid_news.json"))
//		if err != nil {
//			t.Fatalf("Failed to create invalid news file: %v", err)
//		}
//		_, err = newsHandler.GetNewsFromFolder(errorFolder)
//		if err == nil {
//			t.Errorf("Expected an error when retrieving news from a folder with invalid JSON, but got none")
//		}
//	}
//
//	func TestGetNewsSourceFilePath(t *testing.T) {
//		NewsFolder, _ := setupTestData(t)
//		defer cleanupTestData(t, NewsFolder)
//
//		newsHandler := newsFolder{path: NewsFolder}
//		sourceNames := []string{"test-source"}
//
//		got, err := newsHandler.GetNewsSourceFilePath(sourceNames)
//		if err != nil {
//			t.Fatalf("GetNewsSourceFilePath() failed: %v", err)
//		}
//
//		finalFileName := fmt.Sprintf("test-source/%s.json", timeNow)
//		finalFilePath := filepath.Join(NewsFolder, finalFileName)
//		expected := map[string][]string{
//			"test-source": {finalFilePath},
//		}
//
//		if !reflect.DeepEqual(got, expected) {
//			t.Errorf("GetNewsSourceFilePath() = %v, want %v", got, expected)
//		}
//
//		invalidSourceNames := []string{"#$%@+_=!"}
//		paths, err := newsHandler.GetNewsSourceFilePath(invalidSourceNames)
//		if len(paths) != 0 {
//			t.Errorf("Expected an error when retrieving file paths for a non-existent source, but got none")
//		}
//	}
//
// //func setupTestData(t *testing.T) (string, []entity.News) {
// //	t.Helper()
// //
// //	mockTime := time.Date(2024, 7, 3, 12, 0, 0, 0, time.UTC)
// //	news := []entity.News{
// //		{
// //			Title:       "Test News 1",
// //			Description: "This is a test news article 1",
// //			Link:        "https://example.com/news1",
// //			Date:        mockTime,
// //			Source:      "test_source",
// //		},
// //		{
// //			Title:       "Test News 2",
// //			Description: "This is a test news article 2",
// //			Link:        "https://example.com/news2",
// //			Date:        mockTime,
// //			Source:      "test_source",
// //		},
// //	}
// //
// //	NewsFolder := "test-data"
// //	newsHandler := newsFolder{path: NewsFolder}
// //	err := newsHandler.AddNews(news, "test-source")
// //	if err != nil {
// //		t.Fatalf("AddNews() failed: %v", err)
// //	}
// //
// //	return NewsFolder, news
// //}
// //
// //func cleanupTestData(t *testing.T, path string) {
// //	t.Helper()
// //
// //	err := os.RemoveAll(path)
// //	if err != nil {
// //		t.Errorf("Failed to clean up test data folder: %v", err)
// //	}
// //}
package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestAddNews(t *testing.T) {
	NewsFolder, expectedNews := setupTestData(t)
	defer cleanupTestData(t, NewsFolder)

	finalFileName := fmt.Sprintf("test-source/%s.json", time.Now().Format("2006-01-02"))
	finalFilePath := filepath.Join(NewsFolder, finalFileName)

	log.Printf("Final file path: %s\n", finalFilePath)

	fileData, err := os.ReadFile(finalFilePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	var storedNews []entity.News
	err = json.Unmarshal(fileData, &storedNews)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	if !reflect.DeepEqual(storedNews, expectedNews) {
		t.Errorf("Stored news data does not match expected data. Got: %v, Expected: %v", storedNews, expectedNews)
	}

	invalidNewsHandler := newsFolder{path: string([]byte{0x00})}
	err = invalidNewsHandler.AddNews(expectedNews, "invalid-source")
	if err == nil {
		t.Errorf("Expected an error when creating a directory with an invalid path, but got none")
	}

	invalidFilePath := filepath.Join("?'@%=", finalFileName)
	err = os.MkdirAll(filepath.Dir(invalidFilePath), 0755)
	if err == nil {
		t.Fatalf("Expected an error when creating a directory with an invalid path, but got none")
	}

	err = os.WriteFile(invalidFilePath, []byte("test"), 0644)
	if err == nil {
		t.Errorf("Expected an error when writing to a file with an invalid path, but got none")
	}
}

func setupTestData(t *testing.T) (string, []entity.News) {
	t.Helper()

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

	NewsFolder := "test-data"
	newsHandler := newsFolder{path: NewsFolder}
	err := newsHandler.AddNews(news, "test-source")
	if err != nil {
		t.Fatalf("AddNews() failed: %v", err)
	}

	return NewsFolder, news
}

func cleanupTestData(t *testing.T, path string) {
	t.Helper()

	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("Failed to clean up test data folder: %v", err)
	}
}

func TestGetNewsFromFolder(t *testing.T) {
	NewsFolder := "../../internal/testdata"
	newsHandler := newsFolder{path: NewsFolder}
	got, err := newsHandler.GetNewsFromFolder("bbc_news")
	if err != nil {
		t.Fatalf("GetNewsFromFolder() failed: %v", err)
	}

	wants := []entity.News{
		{
			Title:       "Why parents are locking themselves in cells at Korean 'happiness factory'",
			Description: "Some South Koreans are spending time in a cell to try to understand their socially isolated children.",
			Link:        "https://www.bbc.com/news/articles/c2x0le06kn7o",
			Date:        time.Date(2024, 6, 30, 1, 6, 19, 0, time.UTC),
			Source:      "bbc_news",
		},
		{
			Title:       "Watch England fans go wild as Bellingham scores late equaliser",
			Description: "Watch England fans erupt at a fanpark in Wembley as Jude Bellingham scores a late stunner to send the game to extra time against Slovakia in the Euro 2024 last-16 match in Gelsenkirchen.\n\n\n",
			Link:        "https://www.bbc.com/sport/football/videos/cl4yj1ve5z7o",
			Date:        time.Date(2024, 6, 30, 19, 31, 26, 0, time.UTC),
			Source:      "bbc_news",
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

	errorFolder := "test-error-data"
	err = os.MkdirAll(filepath.Join(NewsFolder, errorFolder), 0755)
	if err != nil {
		t.Fatalf("Failed to create test erroneous folder: %v", err)
	}

	_, err = os.Create(filepath.Join(NewsFolder, errorFolder, "invalid_news.json"))
	if err != nil {
		t.Fatalf("Failed to create invalid news file: %v", err)
	}

	_, err = newsHandler.GetNewsFromFolder(errorFolder)
	if err == nil {
		t.Errorf("Expected an error when retrieving news from a folder with invalid JSON, but got none")
	}
}

func TestGetNewsSourceFilePath(t *testing.T) {
	NewsFolder, _ := setupTestData(t)
	defer cleanupTestData(t, NewsFolder)

	newsHandler := newsFolder{path: NewsFolder}
	sourceNames := []string{"test-source"}

	got, err := newsHandler.GetNewsSourceFilePath(sourceNames)
	if err != nil {
		t.Fatalf("GetNewsSourceFilePath() failed: %v", err)
	}

	finalFileName := fmt.Sprintf("test-source/%s.json", time.Now().Format("2006-01-02"))
	finalFilePath := filepath.Join(NewsFolder, finalFileName)
	expected := map[string][]string{
		"test-source": {finalFilePath},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("GetNewsSourceFilePath() = %v, want %v", got, expected)
	}

	invalidSourceNames := []string{"#$%@+_=!"}
	paths, err := newsHandler.GetNewsSourceFilePath(invalidSourceNames)
	if len(paths) != 0 {
		t.Errorf("Expected an error when retrieving file paths for a non-existent source, but got none")
	}
}
