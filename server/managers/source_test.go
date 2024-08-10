package managers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"news-aggregator/internal/entity"
	"os"
	"testing"
)

func TestGetSources(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathToFile: entity.PathToFile("path1")},
		{Name: "source2", PathToFile: entity.PathToFile("path2")},
	}
	writeTestDataToFile(sources)
	s := CreateSourceFolder("test_sources.json")
	result, err := s.GetSources()
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, sources, result, "Expected sources to match")
}

func TestGetSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathToFile: entity.PathToFile("path1")},
		{Name: "source2", PathToFile: entity.PathToFile("path3")},
	}
	writeTestDataToFile(sources)
	s := CreateSourceFolder("test_sources.json")
	result, err := s.GetSource("source1")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, sources[0], result, "Expected source to match")

	result, err = s.GetSource("nonexistent")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "no resources found for name: nonexistent", "Expected specific error message")
	assert.Equal(t, entity.Source{}, result, "Expected empty source")
}

func TestCreateSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()
	s := CreateSourceFolder("test_sources.json")
	result, err := s.CreateSource("source1", "path1")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, entity.Source{Name: "source1", PathToFile: "path1"}, result, "Expected source to match")

	result, err = s.CreateSource("source1", "path2")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "Source with name source1 already exists", "Expected specific error message")
	assert.Equal(t, entity.Source{}, result, "Expected empty source")
}

func TestRemoveSourceByName(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathToFile: entity.PathToFile("path1")},
		{Name: "source2", PathToFile: entity.PathToFile("path3")},
	}
	writeTestDataToFile(sources)
	s := CreateSourceFolder("test_sources.json")
	err := s.RemoveSourceByName("source1")
	assert.Nil(t, err, "Expected no error")

	remainingSources, _ := s.GetSources()
	assert.Equal(t, 1, len(remainingSources), "Expected only one source remaining")
	assert.Equal(t, "source2", string(remainingSources[0].Name), "Expected remaining source to be 'source2'")

	err = s.RemoveSourceByName("nonexistent")
	assert.Nil(t, err, "Expected no error for removing nonexistent source")
}

func TestUpdateSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathToFile: entity.PathToFile("path1")},
		{Name: "source2", PathToFile: entity.PathToFile("path3")},
	}
	writeTestDataToFile(sources)
	s := CreateSourceFolder("test_sources.json")
	err := s.UpdateSource("source1", "newpath")
	assert.Nil(t, err, "Expected no error")

	updatedSources, _ := s.GetSources()
	assert.Equal(t, "newpath", string(updatedSources[0].PathToFile), "Expected updated path")

	err = s.UpdateSource("nonexistent", "newpath")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "source with name nonexistent not found", "Expected specific error message")
}

func TestReadFromFileNonExistent(t *testing.T) {
	path := "non_existent_file.json"
	s := CreateSourceFolder(path)

	sources, err := s.GetSources()
	assert.Nil(t, err, "Expected no error when file does not exist")
	assert.Equal(t, 0, len(sources), "Expected no sources from non-existent file")

	os.Remove(path)
}

func TestWriteToFileError(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()
	os.Chmod("test_sources.json", 0444)

	s := CreateSourceFolder("test_sources.json")
	_, err := s.CreateSource("source1", "path1")
	assert.NotNil(t, err, "Expected an error due to write protection")
}

func TestReadFromFileError(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()
	invalidJSON := []byte(`invalid json`)
	os.WriteFile("test_sources.json", invalidJSON, 0644)

	s := CreateSourceFolder("test_sources.json")
	_, err := s.GetSources()
	assert.NotNil(t, err, "Expected an error due to invalid JSON")
}

func setupTestFile() {
	data := []byte(`[]`)
	pathToResources := "test_sources.json"
	err := os.WriteFile(pathToResources, data, 0644)
	if err != nil {
		log.Fatalf("Error setting up test file: %v", err)
	}
}

func cleanupTestFile() {
	pathToResources := "test_sources.json"
	err := os.Remove(pathToResources)
	if err != nil {
		log.Fatalf("Error cleaning up test file: %v", err)
	}
}

func writeTestDataToFile(sources []entity.Source) {
	jsonData, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	pathToResources := "test_sources.json"
	err = os.WriteFile(pathToResources, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing test data to file: %v", err)
	}
}
