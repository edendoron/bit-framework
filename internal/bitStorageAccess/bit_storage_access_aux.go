package bitStorageAccess

import (
	. "../../configs/rafael.com/bina/bit"
	. "../apiResponseHandlers"
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
}

func writeConfigFailures(w http.ResponseWriter, failureToWrite *string){
	failure := Failure{}
	if err := json.Unmarshal([]byte(*failureToWrite), &failure); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	filename := failure.Description.UnitName + "_" + failure.Description.TestName + "_" + strconv.FormatUint(failure.Description.TestId, 10)
	f, err := os.OpenFile("storage/config/filtering_rules/" + filename + ".txt", os.O_CREATE|os.O_WRONLY, 0644)
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
}

func writeUserGroupFiltering(w http.ResponseWriter, userGroupConfig *string){
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
				if err == nil && info.IsDir() && reportTime.After(startTime) && reportTime.Before(endTime){
					protoReport, err := ioutil.ReadFile(path + "/tests_results.txt")
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

			}
			return nil
		})
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			}
	//response := make([]TestReport, len(reports))
	//for i, report := range reports {
	//	response[i] = testResultToTestReport(report)
	//}
	err = json.NewEncoder(w).Encode(&reports)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
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
}

func readUserGroupMaskedTestIds(w http.ResponseWriter, userGroup string){
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

func testReportToTestResult(tr TestReport) TestResult{
	return TestResult{
		TestId:         uint64(tr.TestId),
		Timestamp:      timestamppb.New(tr.Timestamp),
		TagSet:         convertToKeyValuePair(tr.TagSet),
		FieldSet:       convertToKeyValuePair(tr.FieldSet),
		ReportPriority: uint32(tr.ReportPriority),
	}
}

func testResultToTestReport(tr TestResult) TestReport{
	return TestReport{
		TestId:         float64(tr.TestId),
		ReportPriority: float64(tr.ReportPriority),
		Timestamp:      tr.Timestamp.AsTime(),
		TagSet:         convertToKeyValue(tr.TagSet),
		FieldSet:       convertToKeyValue(tr.FieldSet),
	}
}
