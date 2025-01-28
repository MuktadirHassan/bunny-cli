package bunny

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"sync"
)

// UploadFolder recursively uploads all files in the folder concurrently.
// Example usage: UploadFolder("./path/to/folder", 3)
func UploadFolder(folderPath string, concurrencyLimit int) error {
	slog.Debug("uploading folder: ", slog.String("folder", folderPath))

	var wg sync.WaitGroup
	limiter := make(chan struct{}, concurrencyLimit) // Semaphore for limiting concurrency
	errChan := make(chan error, 1)                   // Channel to capture errors from workers

	// Walk folder
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relativePath, err := filepath.Rel(folderPath, path)
			if err != nil {
				return err
			}
			slog.Debug("queueing file for upload: ", slog.String("path", relativePath))

			wg.Add(1)
			go func() {
				defer wg.Done()
				limiter <- struct{}{} // Acquire semaphore
				defer func() { <-limiter }()

				if err := UploadFileTest(path, relativePath); err != nil {
					select {
					case errChan <- err: // Send error to the error channel
					default: // Avoid blocking if another error has already been sent
					}
				}
			}()
		}
		return nil
	})

	// Wait for all workers to finish
	wg.Wait()

	// Check for errors from WalkDir or workers
	select {
	case workerErr := <-errChan:
		return workerErr
	default:
		return err
	}
}
