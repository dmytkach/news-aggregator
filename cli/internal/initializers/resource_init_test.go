package initializers

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadStaticResourcesFromFolder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "static_resources_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockFiles := map[string][]string{
		"resource1": {"file1.txt", "file2.txt"},
		"resource2": {"file3.txt", "file4.txt", "file5.txt"},
	}

	for resourceName, files := range mockFiles {
		resourceDir := filepath.Join(tempDir, resourceName)
		err := os.MkdirAll(resourceDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", resourceDir, err)
		}
		for _, file := range files {
			filePath := filepath.Join(resourceDir, file)

			f, err := os.Create(filePath)
			if err != nil {
				t.Fatalf("Failed to create file %s: %v", filePath, err)
			}
			defer f.Close()

			_, err = f.WriteString("test content")
			if err != nil {
				t.Fatalf("Failed to write content to file %s: %v", filePath, err)
			}
		}
	}

	result, err := LoadStaticResourcesFromFolder(tempDir)
	if err != nil {
		t.Errorf("LoadStaticResourcesFromFolder returned an error: %v", err)
	}

	expected := map[string][]string{
		"resource1": {
			filepath.Join(tempDir, "resource1/file1.txt"),
			filepath.Join(tempDir, "resource1/file2.txt"),
		},
		"resource2": {
			filepath.Join(tempDir, "resource2/file3.txt"),
			filepath.Join(tempDir, "resource2/file4.txt"),
			filepath.Join(tempDir, "resource2/file5.txt"),
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result does not match expected:\nExpected: %v\nGot: %v", expected, result)
	}

	_, err = LoadStaticResourcesFromFolder("/non-existing-folder")
	if err == nil {
		t.Error("Expected an error for non-existing directory, but got nil")
	}
}
