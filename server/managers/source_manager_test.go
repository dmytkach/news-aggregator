package managers

import (
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"news-aggregator/internal/entity"
	"os"
	"testing"
)

const testResourcesFile = "test_sources.json"

type MockedResourceManager struct {
	mock.Mock
}

func (m *MockedResourceManager) GetSources() ([]entity.Source, error) {
	args := m.Called()
	return args.Get(0).([]entity.Source), args.Error(1)
}

func (m *MockedResourceManager) GetSource(name string) (entity.Source, error) {
	args := m.Called(name)
	return args.Get(0).(entity.Source), args.Error(1)
}

func (m *MockedResourceManager) CreateSource(name, url string) (entity.Source, error) {
	args := m.Called(name, url)
	return args.Get(0).(entity.Source), args.Error(1)
}

func (m *MockedResourceManager) RemoveSourceByName(sourceName string) error {
	args := m.Called(sourceName)
	return args.Error(0)
}

func (m *MockedResourceManager) UpdateSource(oldUrl, newUrl string) error {
	args := m.Called(oldUrl, newUrl)
	return args.Error(0)
}

func TestCRUDOperationsOnSources(t *testing.T) {
	if err := initializeTestFile(); err != nil {
		t.Fatalf("Failed to initialize test file: %v", err)
	}
	defer func() {
		err := os.Remove(testResourcesFile)
		if err != nil {
			t.Fatalf("Failed to clean up test file: %v", err)
		}
	}()

	mockManager := new(MockedResourceManager)

	mockManager.On("CreateSource", "test-source", "test-path").Return(
		entity.Source{Name: "test-source", PathsToFile: []entity.PathToFile{"test-path"}}, nil)
	mockManager.On("GetSource", "test-source").Return(
		entity.Source{Name: "test-source", PathsToFile: []entity.PathToFile{"test-path"}}, nil)
	mockManager.On("UpdateSource", "test-path", "new-test-path").Return(nil)
	mockManager.On("RemoveSourceByName", "test-source").Return(nil)

	_, err := mockManager.CreateSource("test-source", "test-path")
	if err != nil {
		t.Fatalf("CreateSource() failed: %v", err)
	}

	_, err = mockManager.GetSource("test-source")
	if err != nil {
		t.Fatalf("GetSource() failed: %v", err)
	}

	err = mockManager.UpdateSource("test-path", "new-test-path")
	if err != nil {
		t.Fatalf("UpdateSource() failed: %v", err)
	}

	err = mockManager.RemoveSourceByName("test-source")
	if err != nil {
		t.Fatalf("RemoveSourceByName() failed: %v", err)
	}

	mockManager.AssertExpectations(t)
}

func initializeTestFile() error {
	testData := []entity.Source{
		{
			Name:        "test-source",
			PathsToFile: []entity.PathToFile{"test-path"},
		},
	}

	// Write test data to file
	jsonData, err := json.Marshal(testData)
	if err != nil {
		return err
	}
	err = os.WriteFile(testResourcesFile, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
