package bitHistoryCurator

import (
	"net/http"
	"time"
)

func RemoveAgedData(agedTimeDuration time.Duration) {
	//fmt.Println(agedTimeDuration)

	req, err := http.NewRequest(http.MethodDelete, Configs.StorageDeleteURL, nil)
	if err != nil {
		//TODO: handle error
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	params := req.URL.Query()

	timestamp := time.Now().Add(-agedTimeDuration)
	const layout = "2006-January-02 15:4:5"
	timestampStr := timestamp.Format(layout)
	params.Add("timestamp", timestampStr)

	req.URL.RawQuery = params.Encode()

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
