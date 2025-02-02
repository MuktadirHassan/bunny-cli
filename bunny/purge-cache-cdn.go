/*
Purges cache of a BunnyCDN pullzone
Doc: https://docs.bunny.net/reference/pullzonepublic_purgecachepostbytag
API: https://api.bunny.net/pullzone/{id}/purgeCache
*/
package bunny

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// CdnFullCachePurge purges cache of a BunnyCDN pullzone
func CdnFullCachePurge(pullzoneID string) error {
	BUNNYCDN_API_KEY := os.Getenv("BUNNYCDN_API_KEY")
	if BUNNYCDN_API_KEY == "" {
		return fmt.Errorf("env BUNNYCDN_API_KEY is not set")
	}

	purgeURL := fmt.Sprintf("https://api.bunny.net/pullzone/%s/purgeCache", pullzoneID)
	req, err := http.NewRequest("POST", purgeURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create purge request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("AccessKey", BUNNYCDN_API_KEY)

	slog.Debug("purging full cache for pullzone", slog.String("pullzoneID", pullzoneID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to purge cache: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		slog.Debug("successfully purged cache for pullzone", slog.String("pullzoneID", pullzoneID))
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to purge cache: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return nil
}
