package bitHandler

import (
	. "../apiResponseHandlers"
	. "../models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Error creating request for storage access", err)
		return
	}
	//defer req.Body.Close()

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
	r.Body.Close()
}

func PostTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	triggerRequest := TriggerBody{}

	if !handleParam(w, r, triggerRequest) {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&triggerRequest)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad Request", err)
		return
	}

	// validate that data from the user is of type TriggerBody
	err = ValidateType(triggerRequest)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad Request", err)
		return
	}

	// update current TriggerBody
	TriggerChannel <- triggerRequest

	// return ApiResponse response to the user
	ApiResponseHandler(w, http.StatusOK, "Trigger updated!", nil)
	w.WriteHeader(http.StatusOK)
}
