/*
Upload a file to a storage zone based on the URL path. If the directory tree does not exist, it will be created automatically.
Doc: https://docs.bunny.net/reference/put_-storagezonename-path-filename
API: https://{region}.storage.bunnycdn.com/{STORAGE_ZONE_NAME}/{path}/{fileName}
*/
package bunny

import (
	"bytes"
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
	STORAGE_REGION     = os.Getenv("STORAGE_REGION") // Optional; default is "sg"
	STORAGE_ZONE_NAME  = os.Getenv("STORAGE_ZONE_NAME")
	STORAGE_ACCESS_KEY = os.Getenv("STORAGE_ACCESS_KEY")
)

/*
 localPath is the path to the file on the local machine, and remotePath is the path to the file on the BunnyCDN storage zone.
 remotePath will have the same directory structure as localPath
 For example, if localPath is ./path/to/file.txt, then remotePath will be be /path/to/file.txt
*/
// Example usage: UploadFile("./path/to/file.txt", "/path/to/file.txt")
func UploadFile(localPath string, relativePath string) error {
	if STORAGE_ACCESS_KEY == "" || STORAGE_ZONE_NAME == "" {
		return fmt.Errorf("env STORAGE_ACCESS_KEY or STORAGE_ZONE_NAME is not set")
	}

	if STORAGE_REGION == "" {
		STORAGE_REGION = "sg"
	}

	uploadPath := fmt.Sprintf("https://%s.storage.bunnycdn.com/%s/%s", STORAGE_REGION, STORAGE_ZONE_NAME, filepath.ToSlash(relativePath))
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", localPath, err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", localPath, err)
	}

	req, err := http.NewRequest("PUT", uploadPath, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("accept", "application/json")
	req.Header.Set("AccessKey", STORAGE_ACCESS_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %v", localPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		slog.Info("successfully uploaded file", slog.String("localPath", localPath), slog.String("remotePath", relativePath))
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload %s: %d - %s", localPath, resp.StatusCode, string(respBody))
	}
	return nil
}

func UploadFileTest(localPath string, relativePath string) error {
	time.Sleep(200 * time.Millisecond) // Simulate network latency
	// if count == 5 {
	// 	return fmt.Errorf("failed to upload %s", localPath)
	// }
	// count++
	slog.Info("successfully uploaded file", slog.String("path", localPath))
	return nil
}
