package bitStorageAccess

import (
	. "../../configs/rafael.com/bina/bit"
	. "../apiResponseHandlers"
	. "../bitHandler"
	. "../models"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func writeReports(w http.ResponseWriter, testReports *string) {
	reports := ReportBody{}
	if err := json.Unmarshal([]byte(*testReports), &reports); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	for _, report := range reports.Reports {
		requestToProto := testReportToTestResult(report)
		protoReports, err := proto.Marshal(&requestToProto)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		path := "storage/test_reports/" + fmt.Sprint(report.Timestamp.Date()) + "/" + fmt.Sprint(report.Timestamp.Hour()) +
			"/" + fmt.Sprint(report.Timestamp.Minute()) + "/" + fmt.Sprint(report.Timestamp.Second())
		if _, err = os.Stat(path + "/tests_results.txt"); os.IsNotExist(err) {
			err = os.MkdirAll(path, 0700)
			if err != nil {
				ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			}
		}
		f, err := os.OpenFile(path+"/tests_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		if _, err = f.Write(protoReports); err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func writeConfigFailures(w http.ResponseWriter, failureToWrite *string) {
	failure := Failure{}
	if err := json.Unmarshal([]byte(*failureToWrite), &failure); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	filename := failure.Description.UnitName + "_" + failure.Description.TestName + "_" + strconv.FormatUint(failure.Description.TestId, 10)
	f, err := os.OpenFile("storage/config/filtering_rules/"+filename+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	failureToProto, err := proto.Marshal(&failure)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(failureToProto); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func writeExtendedFailures(w http.ResponseWriter, failureToWrite *string) {
	failure := ExtendedFailure{}
	if err := json.Unmarshal([]byte(*failureToWrite), &failure); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	filename := failure.Failure.Description.UnitName + "_" + failure.Failure.Description.TestName + "_" + strconv.FormatUint(failure.Failure.Description.TestId, 10)
	f, err := os.OpenFile("storage/config/perm_filtering_rules/"+filename+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	failureToProto, err := proto.Marshal(&failure)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(failureToProto); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func writeUserGroupFiltering(w http.ResponseWriter, userGroupConfig *string) {
	userGroupsFilters := UserGroupsFiltering_FilteredFailures{}
	if err := json.Unmarshal([]byte(*userGroupConfig), &userGroupsFilters); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		if userGroupsFilters.UserGroup == "" {
			return
		}
	}
	filename := "storage/config/user_groups_masks/" + userGroupsFilters.UserGroup + ".txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	userGroupToProto, err := proto.Marshal(&userGroupsFilters)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(userGroupToProto); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func writeBitStatus(w http.ResponseWriter, bitStatus *string) {
	status := BitStatus{}
	if err := json.Unmarshal([]byte(*bitStatus), &status); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	protoStatus, err := proto.Marshal(&status)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	currentTime := time.Now()
	path := "storage/bit_status/" + fmt.Sprint(currentTime.Date()) + "/" + fmt.Sprint(currentTime.Hour()) +
		"/" + fmt.Sprint(currentTime.Minute()) + "/" + fmt.Sprint(currentTime.Second())
	if _, err = os.Stat(path + "/bit_status.txt"); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
	}
	f, err := os.OpenFile(path+"/bit_status.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(protoStatus); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func readReports(w http.ResponseWriter, start string, end string, filter string) {
	var reports []TestResult
	const layout = "2006-January-02 15:4:5"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = filepath.Walk("storage/test_reports/",
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && reportTime.After(startTime) && reportTime.Before(endTime) {
					protoReport, err := ioutil.ReadFile(path + "/tests_results.txt")
					if err != nil {
						ApiResponseHandler(w, http.StatusInternalServerError, "Can't find report!", err)
					}
					var temp TestResult
					err = proto.Unmarshal(protoReport, &temp)
					if err != nil {
						ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
					reports = append(reports, temp)
				}
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
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func readConfigFailures(w http.ResponseWriter) {
	var configFailures []Failure
	files, err := ioutil.ReadDir("storage/config/filtering_rules")
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	for _, f := range files {
		content, err := ioutil.ReadFile("storage/config/filtering_rules/" + f.Name())
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		decodedContent := Failure{}
		err = proto.Unmarshal(content, &decodedContent)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		configFailures = append(configFailures, decodedContent)
	}
	err = json.NewEncoder(w).Encode(&configFailures)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func readExtendedFailures(w http.ResponseWriter) {
	var foreverFailures []ExtendedFailure
	files, err := ioutil.ReadDir("storage/config/perm_filtering_rules")
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	for _, f := range files {
		content, err := ioutil.ReadFile("storage/config/perm_filtering_rules/" + f.Name())
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		decodedContent := ExtendedFailure{}
		err = proto.Unmarshal(content, &decodedContent)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
		foreverFailures = append(foreverFailures, decodedContent)
	}
	err = json.NewEncoder(w).Encode(&foreverFailures)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func readBitStatus(w http.ResponseWriter, start string, end string, filter string) {
	var statuses []BitStatus
	const layout = "2006-January-02 15:4:5"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = filepath.Walk("storage/bit_status/",
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && reportTime.After(startTime) && reportTime.Before(endTime) {
					protoStatus, err := ioutil.ReadFile(path + "/bit_status.txt")
					if err != nil {
						ApiResponseHandler(w, http.StatusInternalServerError, "Can't find status!", err)
					}
					temp := BitStatus{}
					err = proto.Unmarshal(protoStatus, &temp)
					if err != nil {
						ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
					statuses = append(statuses, temp)
				}
			}
			return nil
		})
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = json.NewEncoder(w).Encode(&statuses)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func readUserGroupMaskedTestIds(w http.ResponseWriter, userGroup string) {
	content, err := ioutil.ReadFile("storage/config/user_groups_masks/" + userGroup + ".txt")
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	decodedContent := UserGroupsFiltering_FilteredFailures{}
	err = proto.Unmarshal(content, &decodedContent)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = json.NewEncoder(w).Encode(&decodedContent.MaskedTestIds)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	w.WriteHeader(http.StatusOK)
}

func convertToKeyValuePair(arr []KeyValue) []*KeyValuePair {
	copyArr := make([]*KeyValuePair, len(arr))
	for k, v := range arr {
		pair := KeyValuePair{Key: []byte(v.Key), Value: []byte(v.Value)}
		copyArr[k] = &pair
	}
	return copyArr
}

func convertToKeyValue(arr []*KeyValuePair) []KeyValue {
	copyArr := make([]KeyValue, len(arr))
	for k, v := range arr {
		pair := KeyValue{Key: string(v.Key), Value: string(v.Value)}
		copyArr[k] = pair
	}
	return copyArr
}

func testReportToTestResult(tr TestReport) TestResult {
	return TestResult{
		TestId:         uint64(tr.TestId),
		Timestamp:      timestamppb.New(tr.Timestamp),
		TagSet:         convertToKeyValuePair(tr.TagSet),
		FieldSet:       convertToKeyValuePair(tr.FieldSet),
		ReportPriority: uint32(tr.ReportPriority),
	}
}

func testResultToTestReport(tr TestResult) TestReport {
	return TestReport{
		TestId:         float64(tr.TestId),
		ReportPriority: float64(tr.ReportPriority),
		Timestamp:      tr.Timestamp.AsTime(),
		TagSet:         convertToKeyValue(tr.TagSet),
		FieldSet:       convertToKeyValue(tr.FieldSet),
	}
}
