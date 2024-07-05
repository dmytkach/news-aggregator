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
		{Name: "source1", PathsToFile: []entity.PathToFile{"path1", "path2"}},
		{Name: "source2", PathsToFile: []entity.PathToFile{"path3"}},
	}
	writeTestDataToFile(sources)

	result, err := GetSources()
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, sources, result, "Expected sources to match")
}

func TestGetSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathsToFile: []entity.PathToFile{"path1", "path2"}},
		{Name: "source2", PathsToFile: []entity.PathToFile{"path3"}},
	}
	writeTestDataToFile(sources)

	result, err := GetSource("source1")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, sources[0], result, "Expected source to match")

	result, err = GetSource("nonexistent")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "No resources found for name: nonexistent", "Expected specific error message")
	assert.Equal(t, entity.Source{}, result, "Expected empty source")
}

func TestCreateSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	result, err := CreateSource("source1", "path1")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, entity.Source{Name: "source1", PathsToFile: []entity.PathToFile{"path1"}}, result, "Expected source to match")

	result, err = CreateSource("source1", "path2")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, entity.Source{Name: "source1", PathsToFile: []entity.PathToFile{"path1", "path2"}}, result, "Expected empty source")
	result, err = CreateSource("source1", "path1")

	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "resource already exists", "Expected specific error message")
	assert.Equal(t, entity.Source{}, result, "Expected empty source")
}

func TestRemoveSourceByName(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathsToFile: []entity.PathToFile{"path1", "path2"}},
		{Name: "source2", PathsToFile: []entity.PathToFile{"path3"}},
	}
	writeTestDataToFile(sources)

	err := RemoveSourceByName("source1")
	assert.Nil(t, err, "Expected no error")

	remainingSources, _ := GetSources()
	assert.Equal(t, 1, len(remainingSources), "Expected only one source remaining")
	assert.Equal(t, "source2", string(remainingSources[0].Name), "Expected remaining source to be 'source2'")

}

func TestUpdateSource(t *testing.T) {
	setupTestFile()
	defer cleanupTestFile()

	sources := []entity.Source{
		{Name: "source1", PathsToFile: []entity.PathToFile{"path1", "path2"}},
		{Name: "source2", PathsToFile: []entity.PathToFile{"path3"}},
	}
	writeTestDataToFile(sources)

	err := UpdateSource("path2", "newpath")
	assert.Nil(t, err, "Expected no error")

	updatedSources, _ := GetSources()
	assert.Equal(t, "newpath", string(updatedSources[0].PathsToFile[1]), "Expected updated path")

	err = UpdateSource("nonexistent", "newpath")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "source with URL nonexistent not found", "Expected specific error message")

	err = UpdateSource("newpath", "path1")
	assert.NotNil(t, err, "Expected an error")
	assert.EqualError(t, err, "resource already exists", "Expected specific error message")
}
func setupTestFile() {
	data := []byte(`[]`)
	PathToResources = "test_sources.json"
	err := os.WriteFile(PathToResources, data, 0644)
	if err != nil {
		log.Fatalf("Error setting up test file: %v", err)
	}
}

func cleanupTestFile() {
	PathToResources = "test_sources.json"
	err := os.Remove(PathToResources)
	if err != nil {
		log.Fatalf("Error cleaning up test file: %v", err)
	}
}
func writeTestDataToFile(sources []entity.Source) {
	jsonData, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	PathToResources = "test_sources.json"
	err = os.WriteFile(PathToResources, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing test data to file: %v", err)
	}
}
