package bitStorageAccess

import (
	. "../../configs/rafael.com/bina/bit"
	. "../apiResponseHandlers"
	. "../models"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetDataRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	query := r.URL.Query()
	if len(query["reports"]) > 0 {
		readReports(w, query["start"][0], query["end"][0], query["filter"][0])
	} else if len(query["config_failures"]) > 0 {
		readConfigFailures(w)
	} else if len(query["forever_failure"]) > 0 {
		readExtendedFailures(w)
	} else if len(query["config_user_groups_filtering"]) > 0 {
		readUserGroupMaskedTestIds(w, query["id"][0])
	} else if len(query["bit_status"]) > 0 {
		readBitStatus(w, query["start"][0], query["end"][0], query["filter"][0])
	} else if len(query["user_groups"]) > 0 {
		readUserGroups(w)
	}
}

func PostDataWrite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	requestBody := KeyValue{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	switch requestBody.Key {
	case "reports":
		writeReports(w, &requestBody.Value)
	case "config_failure":
		writeConfigFailures(w, &requestBody.Value)
	case "forever_failure":
		writeExtendedFailures(w, &requestBody.Value)
		ApiResponseHandler(w, http.StatusOK, "Failure received!", nil)
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
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	var reports []TestResult
	err = filepath.Walk("storage/"+query+"/",
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				protoReport, err := ioutil.ReadFile(path)
				if err != nil {
					ApiResponseHandler(w, http.StatusInternalServerError, "Can't find report!", err)
				}
				temp := TestResult{}
				err = proto.Unmarshal(protoReport, &temp)
				if err != nil {
					ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				}
				reports = append(reports, temp)
			}
			return nil
		})
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}

	response := make([]TestReport, len(reports))
	for i, report := range reports {
		response[i] = testResultToTestReport(report)
	}
}

func PutDataWrite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	request := TestReport{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	path := "../storage/" + fmt.Sprint(request.Timestamp.Date()) + "/" + fmt.Sprint(request.Timestamp.Hour()) +
		"/" + fmt.Sprint(request.Timestamp.Minute()) + "/" + fmt.Sprint(request.Timestamp.Second())
	if _, err = os.Stat(path + "/tests_results.txt"); os.IsNotExist(err) {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	f, err := os.OpenFile(path+"/tests_results.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	requestToProto := testReportToTestResult(request)
	protoReports, err := proto.Marshal(&requestToProto)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(protoReports); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	timestamp := r.URL.Query()["timestamp"][0]
	const layout = "2006-January-02 15:4:5"
	threshold, err := time.Parse(layout, timestamp)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = filepath.Walk("storage/test_reports",
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && reportTime.Before(threshold) {
					err = os.RemoveAll(path)
					if err != nil {
						ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
				}
			}
			return nil
		})
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	deleteEmptyDirectories(w, "storage/test_reports", 0)
	ApiResponseHandler(w, http.StatusOK, "Old reports deleted successfully", nil)
}

func deleteEmptyDirectories(w http.ResponseWriter, path string, depth int) {
	if depth == 3 {
		return
	}
	err := filepath.Walk(path,
		func(childPath string, info fs.FileInfo, err error) error {
			deleteEmptyDirectories(w, childPath, depth+1)
			if path != "storage/test_reports" && IsEmpty(path) {
				err := os.RemoveAll(path)
				if err != nil {
					ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				}
			}
			return nil
		})
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
}

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}
