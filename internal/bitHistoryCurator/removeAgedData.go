package bitHistoryCurator

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const bitStorageAccessUrl = "http://localhost:8082"

func removeAgedData(agedTime time.Duration) {
	timeStamp := time.Now().Add(-agedTime)
	fmt.Println(timeStamp)

	params := url.Values{}
	params.Add("timestamp", timeStamp.String())

	req, err := http.NewRequest(http.MethodDelete, bitStorageAccessUrl, strings.NewReader(params.Encode()))
	if err != nil {
		//TODO: handle error
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		time.Sleep(10 * time.Second)
		//TODO: handle error
		return
	}
	err = resp.Body.Close()
	if err != nil {
		//TODO: handle error
		return
	}
}
