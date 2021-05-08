package bitHandler

import (
	. "../apiResponseHandlers"
	"io/ioutil"
	"net/http"
)

const storageDataGetTriggerURL = "http://localhost:8082/data/read"

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	req, err := http.NewRequest(http.MethodGet, storageDataGetTriggerURL, nil)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Error creating request for storage access", err)
		return
	}
	defer req.Body.Close()

	params := req.URL.Query()
	params.Add("key", "BITStatus")

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Storage Access Error", err)
		return
	}
	defer storageResponse.Body.Close()

	respBody, err := ioutil.ReadAll(storageResponse.Body)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Storage Access Error", err)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
