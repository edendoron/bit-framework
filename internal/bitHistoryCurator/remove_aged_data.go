package bithistorycurator

import (
	"log"
	"net/http"
	"time"
)

// RemoveAgedData sends a delete request to storage to delete obsolete data.
func RemoveAgedData(agedTimeDuration time.Duration) {
	req, err := http.NewRequest(http.MethodDelete, Configs.StorageDeleteURL, nil)
	if err != nil {
		log.Printf("error creating delete request to storage: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	params := req.URL.Query()

	const layout = "2006-January-02 15:4:5"
	timestamp := time.Now().Add(-agedTimeDuration).Format(layout)
	params.Add("timestamp", timestamp)

	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("error deleting from storage: %v", err)
		return
	}
	resp.Body.Close()
}
