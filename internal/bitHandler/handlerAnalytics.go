package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

type Analyzer struct {
	Failures []Failure
	Reports  []TestResult
	Status   BitStatus
}

func (a *Analyzer) ReadFromStorage(key string) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		// TODO: handle error
		return
	}
	//defer req.Body.Close()

	params := req.URL.Query()
	params.Add("key", key)

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}
	defer storageResponse.Body.Close()

	if key == "Failure" {
		err = json.NewDecoder(storageResponse.Body).Decode(&a.Failures)
	} else if key == "Report" {
		err = json.NewDecoder(storageResponse.Body).Decode(&a.Reports)
	}

	if err != nil {
		// TODO: handle error
		return
	}
}

func (a *Analyzer) Crosscheck() {
	for _, failure := range a.Failures {
		countFailed, timestamp := a.checkExaminationRule(failure)
		if countFailed > 0 {
			reportedFailure := &BitStatus_RportedFailure{
				FailureData: failure.Description,
				Timestamp:   timestamp,
				Count:       countFailed,
			}
			a.Status.Failures = append(a.Status.Failures, reportedFailure)
		}
		// TODO: a.Status.UserGroup
	}
}

func (a *Analyzer) WriteBitStatus() {

	//TODO: handle error
	jsonStatus, _ := json.MarshalIndent(a.Status, "", " ")

	message := KeyValuePair{
		Key:   []byte("BitStatus"),
		Value: jsonStatus,
	}

	//TODO: handle error
	jsonMessage, _ := json.MarshalIndent(message, "", " ")
	postBody := bytes.NewReader(jsonMessage)

	storageResponse, err := http.Post(storageDataWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		//TODO: handle this error
		return
	}
	defer storageResponse.Body.Close()
}

func (a *Analyzer) checkExaminationRule(failure Failure) (uint64, *timestamppb.Timestamp) {

}
