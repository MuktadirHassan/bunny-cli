package bunny

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"
	"sync"
	"time"
)

// UploadFolder recursively uploads all files in the folder concurrently using a worker pool.
// Example usage: UploadFolder("./path/to/folder", 3, 10*time.Second, true)
func UploadFolder(folderPath string, concurrencyLimit int, timeout time.Duration, failFast bool) error {
	slog.Debug("uploading folder: ", slog.String("folder", folderPath))

	var wg sync.WaitGroup
	fileChan := make(chan string)  // Channel to send file paths to workers
	errChan := make(chan error, 1) // Channel to capture errors from workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is canceled when the function exits

	// Count total files for the progress bar
	totalFiles := 0
	filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})

	slog.Info("total files", slog.Int("files", totalFiles))

	// Start worker pool
	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case path, ok := <-fileChan:
					if !ok {
						return // Channel closed, exit the worker
					}
					relativePath, err := filepath.Rel(folderPath, path)
					if err != nil {
						select {
						case errChan <- err:
						default:
						}
						continue
					}

					var uploadErr error

					// Retry loop with delays and timeout
					for retry := 0; retry < 3; retry++ {
						uploadCtx, uploadCancel := context.WithTimeout(ctx, timeout)
						defer uploadCancel()

						// Attempt the upload with timeout
						if err := UploadFileTest(uploadCtx, path, relativePath); err == nil {
							break // Upload succeeded, exit retry loop
						} else if err == context.Canceled && failFast {
							uploadErr = err
							slog.Warn("upload canceled", slog.String("path", path), slog.Int("retry", retry+1), slog.String("error", err.Error()))
							return // Exit the worker immediately
						} else {
							uploadErr = err
							slog.Warn("upload failed, retrying", slog.String("path", path), slog.Int("retry", retry+1), slog.String("error", err.Error()))

							// Add a delay between retries (e.g., 2 seconds)
							if retry < 2 { // No delay needed after the last retry
								time.Sleep(2 * time.Second)
							}
						}
					}

					if uploadErr != nil {
						select {
						case errChan <- uploadErr:
							if failFast {
								cancel() // Cancel all ongoing uploads
								return   // Exit the worker immediately
							}
						default:
						}
					}

				case <-ctx.Done():
					return // Context canceled, exit the worker
				}
			}
		}()
	}

	// Walk folder and send file paths to the worker pool
	go func() {
		defer close(fileChan) // Close the channel when done
		err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				select {
				case fileChan <- path:
				case <-ctx.Done(): // Stop sending files if context is canceled
					return ctx.Err()
				}
			}
			return nil
		})
		if err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	// Wait for all workers to finish
	wg.Wait()

	// Check for errors from WalkDir or workers
	select {
	case workerErr := <-errChan:
		return workerErr
	default:
		return nil
	}
}
