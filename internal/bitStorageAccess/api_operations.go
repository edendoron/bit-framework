package bitStorageAccess

import (
	"encoding/json"
	"fmt"
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// GetDataRead chooses the correct helper function to call based on the received parameters.
func GetDataRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	query := r.URL.Query()
	if len(query["reports"]) > 0 {
		readReports(w, query["start"][0], query["end"][0])
	} else if len(query["config_failures"]) > 0 {
		readConfigFailures(w)
	} else if len(query["forever_failure"]) > 0 {
		readExtendedFailures(w)
	} else if len(query["config_user_groups_filtering"]) > 0 {
		readUserGroupMaskedTestIds(w, query["id"][0])
	} else if len(query["bit_status"]) > 0 {
		readBitStatus(w, query["start"][0], query["end"][0])
	} else if len(query["user_groups"]) > 0 {
		readUserGroups(w)
	}
}

// PostDataWrite chooses the correct helper function to call based on the received parameters.
func PostDataWrite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	requestBody := models.KeyValue{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	switch requestBody.Key {
	case "reports":
		writeReports(w, &requestBody.Value)
	case "config_failure":
		writeConfigFailures(w, &requestBody.Value)
	case "forever_failure":
		writeExtendedFailures(w, &requestBody.Value)
		rh.ApiResponseHandler(w, http.StatusOK, "Failure received!", nil)
	case "config_user_group_filtering":
		writeUserGroupFiltering(w, &requestBody.Value)
	case "bit_status":
		writeBitStatus(w, &requestBody.Value)
	}

}

func PutDataRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var query string
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	var reports []bit.TestResult
	err = filepath.Walk("storage/"+query+"/",
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				protoReport, err := ioutil.ReadFile(path)
				if err != nil {
					rh.ApiResponseHandler(w, http.StatusInternalServerError, "Can't find report!", err)
				}
				temp := bit.TestResult{}
				err = proto.Unmarshal(protoReport, &temp)
				if err != nil {
					rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				}
				reports = append(reports, temp)
			}
			return nil
		})
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}

	response := make([]models.TestReport, len(reports))
	for i, report := range reports {
		response[i] = testResultToTestReport(report)
	}
}

func PutDataWrite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	request := models.TestReport{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	path := "../storage/" + fmt.Sprint(request.Timestamp.Date()) + "/" + fmt.Sprint(request.Timestamp.Hour()) +
		"/" + fmt.Sprint(request.Timestamp.Minute()) + "/" + fmt.Sprint(request.Timestamp.Second())
	if _, err = os.Stat(path + "/tests_results.txt"); os.IsNotExist(err) {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	f, err := os.OpenFile(path+"/tests_results.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	requestToProto := testReportToTestResult(request)
	protoReports, err := proto.Marshal(&requestToProto)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(protoReports); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	rh.ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}

// DeleteData receives a timestamp and calls the deletion function for both test reports and bit statuses.
func DeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	timestamp := r.URL.Query()["timestamp"][0]
	deleteAgedData(w, "test_reports", timestamp)
	deleteAgedData(w, "bit_status", timestamp)
	rh.ApiResponseHandler(w, http.StatusOK, "Outdated data deleted successfully", nil)
}
