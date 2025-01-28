/*
Upload a file to a storage zone based on the URL path. If the directory tree does not exist, it will be created automatically.
Doc: https://docs.bunny.net/reference/put_-storagezonename-path-filename
API: https://{region}.storage.bunnycdn.com/{STORAGE_ZONE_NAME}/{path}/{fileName}
*/
package bunny

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// Validate env vars
	STORAGE_ZONE_NAME     = os.Getenv("STORAGE_ZONE_NAME")
	STORAGE_ACCESS_KEY    = os.Getenv("STORAGE_ACCESS_KEY")
	STORAGE_ZONE_HOSTNAME = os.Getenv("STORAGE_ZONE_HOSTNAME")
)

// UploadFile uploads a file to the BunnyCDN storage with context support for timeouts and cancellation.
func UploadFile(ctx context.Context, localPath string, relativePath string, STORAGE_ZONE_HOST string) error {
	uploadPath := fmt.Sprintf("https://%s/%s/%s", STORAGE_ZONE_HOST, STORAGE_ZONE_NAME, filepath.ToSlash(relativePath))

	// Open the file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", localPath, err)
	}
	defer file.Close()

	// Read the file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", localPath, err)
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "PUT", uploadPath, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("accept", "application/json")
	req.Header.Set("AccessKey", STORAGE_ACCESS_KEY)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %v", localPath, err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusCreated {
		slog.Info("successfully uploaded file", slog.String("localPath", localPath))
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload %s: %d - %s", localPath, resp.StatusCode, string(respBody))
	}

	return nil
}

// UploadFileTest simulates uploading a file with context support (replace with actual implementation).
func UploadFileTest(ctx context.Context, path string, relativePath string, STORAGE_ZONE_HOST string) error {
	// Simulate file upload by sleeping for a short duration
	select {
	case <-time.After(3 * time.Second): // Simulate upload delay
	case <-ctx.Done(): // Respect context cancellation or timeout
		return ctx.Err()
	}

	// Simulate success or failure based on the file name
	if filepath.Base(path) == "fail.txt" {
		slog.Warn("simulated upload failure", slog.String("path", path))
		return fmt.Errorf("failed to upload file %s", path)
	}

	slog.Info("successfully uploaded file", slog.String("path", path))
	return nil
}
