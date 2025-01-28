package bunny

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// MockUploadFileTest is a mock implementation of UploadFileTest for testing.
func MockUploadFileTest(path string, relativePath string) error {
	// Simulate a successful upload for most files
	if filepath.Base(path) == "fail.txt" {
		return errors.New("upload failed")
	}
	return nil
}

func TestUploadFolder(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create a test directory structure
	testFiles := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
		"subdir/fail.txt", // This file will fail to upload
	}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test cases
	tests := []struct {
		name             string
		folderPath       string
		concurrencyLimit int
		expectError      bool
	}{
		{
			name:             "successful upload",
			folderPath:       tempDir,
			concurrencyLimit: 3,
			expectError:      false,
		},
		{
			name:             "non-existent folder",
			folderPath:       "/nonexistent/folder",
			concurrencyLimit: 3,
			expectError:      true,
		},
		{
			name:             "concurrency limit 1",
			folderPath:       tempDir,
			concurrencyLimit: 1,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UploadFolder(tt.folderPath, tt.concurrencyLimit)
			if tt.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
