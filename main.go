package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	STORAGE_ZONE_NAME  = os.Getenv("STORAGE_ZONE_NAME")
	STORAGE_ACCESS_KEY = os.Getenv("STORAGE_ACCESS_KEY")
	localPath          = flag.String("d", "", "Path to the directory to upload")
	workerCountFlag    = flag.Int("n", 5, "Number of concurrent upload workers")
	workerCount        = *workerCountFlag
)

func main() {
	flag.Parse()
	if *localPath == "" {
		fmt.Println("Please provide a directory path to upload")
		flag.Usage()
		return
	}
	if STORAGE_ACCESS_KEY == "" {
		fmt.Println("env STORAGE_ACCESS_KEY is not set")
		return
	}
	if STORAGE_ZONE_NAME == "" {
		fmt.Println("env STORAGE_ZONE_NAME is not set")
		return
	}

	fmt.Println("Storage Zone Name: ", STORAGE_ZONE_NAME)
	fmt.Println("Starting BunnyCDN uploader...")

	err := uploadFolder(*localPath)

	if err != nil {
		fmt.Println("Failed to upload folder: ", err)
		return
	}

	fmt.Println("Completed running uploader.")
}

// localPath is the path to the file on the local machine, and remotePath is the path to the file on the BunnyCDN storage zone.
// remotePath will have the same directory structure as localPath
func uploadFile(localPath string, relativePath string) error {
	// Upload file to BunnyCDN a correctly formatted upload URL should resemble: https://{region}.bunnycdn.com/{STORAGE_ZONE_NAME}/{path}/{fileName}
	uploadPath := fmt.Sprintf("https://sg.storage.bunnycdn.com/%s/%s", STORAGE_ZONE_NAME, filepath.ToSlash(relativePath))
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
		return fmt.Errorf("upload failed: %v", err)
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
		fmt.Printf("Upload successful for: %s\n", localPath)
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload %s: %d - %s", localPath, resp.StatusCode, string(respBody))
	}
	return nil
}

// Worker pool for concurrent uploads
func worker(fileChan <-chan [2]string, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	fmt.Printf("Worker %d started.\n", workerID)
	for file := range fileChan {
		fmt.Printf("Worker %d uploading file: %s\n", workerID, file[1])
		err := uploadFile(file[0], file[1])
		if err != nil {
			fmt.Printf("Worker %d failed to upload %s: %v\n", workerID, file[1], err)
		}

	}
	fmt.Printf("Worker %d finished.\n", workerID)
}

func uploadFolder(folderPath string) error {
	fmt.Println("Uploading folder: ", folderPath)
	// Create a channel to send files to workers
	fileChan := make(chan [2]string, workerCount)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(fileChan, &wg, i+1)
	}

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relativePath, err := filepath.Rel(folderPath, path)
			if err != nil {
				return err
			}
			fmt.Printf("Queueing file for upload: %s\n", relativePath)
			fileChan <- [2]string{path, relativePath}
		}
		return nil
	})

	// Close the channel and wait for workers to finish
	close(fileChan)
	wg.Wait()
	if err != nil {
		return fmt.Errorf("error walking directory: %v", err)
	}
	return nil
}
