package bitHistoryCurator

import (
	"fmt"
	"net/http"
	"time"
)

const bitStorageAccessUrl = "http://localhost:8082"

func RemoveAgedData(agedTime time.Duration) {
	timeStamp := time.Now().Add(-agedTime)
	fmt.Println(timeStamp)

	req, err := http.NewRequest(http.MethodDelete, bitStorageAccessUrl, nil)
	if err != nil {
		//TODO: handle error
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	params := req.URL.Query()
	params.Add("timestamp", timeStamp.String())

	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	//fmt.Println(req.URL.String())

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
