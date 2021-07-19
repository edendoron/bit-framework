package bitstorageaccess

import (
	"encoding/json"
	"fmt"
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	handler "github.com/edendoron/bit-framework/internal/bitHandler"
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func writeReports(w http.ResponseWriter, testReports *string) {
	reports := models.ReportBody{}
	if err := json.Unmarshal([]byte(*testReports), &reports); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	for _, report := range reports.Reports {
		if time.Now().Sub(report.Timestamp).Seconds() >= Configs.BitHandlerTriggerPeriod {
			log.Printf("Report %v: Timestamp is too late and bitHandler will ignore this report!", report.TestId)
		}
		reportToWrite := testReportToTestResult(report)
		path := "storage/test_reports/" + fmt.Sprint(report.Timestamp.Date()) + "/" + fmt.Sprint(report.Timestamp.Hour()) +
			"/" + fmt.Sprint(report.Timestamp.Minute()) + "/" + fmt.Sprint(report.Timestamp.Second())
		if _, err := os.Stat(path + "/tests_results.txt"); os.IsNotExist(err) {
			err = os.MkdirAll(path, 0700)
			if err != nil {
				rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				return
			}
		}
		f, err := os.OpenFile(path+"/tests_results.txt", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		protoReports, err := ioutil.ReadFile(path + "/tests_results.txt")
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Can't find report!", err)
			return
		}
		var temp bit.TestResultsSet
		err = proto.Unmarshal(protoReports, &temp)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		temp.ResultsSet = append(temp.ResultsSet, &reportToWrite)
		protoReports, err = proto.Marshal(&temp)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		if _, err = f.Write(protoReports); err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		if err = f.Close(); err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func writeConfigFailures(w http.ResponseWriter, failureToWrite *string) {
	failure := bit.Failure{}
	if err := json.Unmarshal([]byte(*failureToWrite), &failure); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	filename := failure.Description.UnitName + "_" + failure.Description.TestName + "_" + strconv.FormatUint(failure.Description.TestId, 10)
	f, err := os.OpenFile("storage/config/filtering_rules/"+filename+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	failureToProto, err := proto.Marshal(&failure)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	if _, err = f.Write(failureToProto); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)

	writeUserGroups(w, failure)
}

func writeUserGroups(w http.ResponseWriter, failure bit.Failure) {
	var content []byte
	if _, err := os.Stat("storage/config/user_groups/user_groups.txt"); err == nil {
		content, err = ioutil.ReadFile("storage/config/user_groups/user_groups.txt")
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
	}

	userGroups := make(map[string]int)
	var groups []string
	if len(content) > 0 {
		groups = strings.Split(string(content), "\n")
	}
	for _, group := range groups {
		userGroups[group] = 1
	}
	for _, group := range failure.Dependencies.BelongsToGroup {
		userGroups[group] = 1
	}
	for _, group := range failure.Dependencies.MasksOtherGroup {
		userGroups[group] = 1
	}
	f, err := os.OpenFile("storage/config/user_groups/user_groups.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	i := 0
	for group := range userGroups {
		if _, err = f.WriteString(group); err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		if i != len(userGroups)-1 {
			_, err = f.WriteString("\n")
			if err != nil {
				rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				return
			}
		}
		i++
	}
}

func writeExtendedFailures(w http.ResponseWriter, failureToWrite *string) {
	failure := handler.ExtendedFailure{}
	if err := json.Unmarshal([]byte(*failureToWrite), &failure); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	filename := failure.Failure.Description.UnitName + "_" + failure.Failure.Description.TestName + "_" + strconv.FormatUint(failure.Failure.Description.TestId, 10)
	f, err := os.OpenFile("storage/config/perm_filtering_rules/"+filename+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	//failureToProto, err := proto.Marshal(&failure)
	//if err != nil {
	//	rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	//}
	if _, err = f.Write([]byte(*failureToWrite)); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeUserGroupFiltering(w http.ResponseWriter, userGroupConfig *string) {
	userGroupsFilters := bit.UserGroupsFiltering_FilteredFailures{}
	if err := json.Unmarshal([]byte(*userGroupConfig), &userGroupsFilters); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	filename := "storage/config/user_groups_masks/" + userGroupsFilters.UserGroup + ".txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	userGroupToProto, err := proto.Marshal(&userGroupsFilters)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	if _, err = f.Write(userGroupToProto); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeBitStatus(w http.ResponseWriter, bitStatus *string) {
	status := bit.BitStatus{}
	if err := json.Unmarshal([]byte(*bitStatus), &status); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	protoStatus, err := proto.Marshal(&status)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	currentTime := time.Now()
	path := "storage/bit_status/" + fmt.Sprint(currentTime.Date()) + "/" + fmt.Sprint(currentTime.Hour()) +
		"/" + fmt.Sprint(currentTime.Minute()) + "/" + fmt.Sprint(currentTime.Second())
	if _, err = os.Stat(path + "/bit_status.txt"); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
	}
	f, err := os.OpenFile(path+"/bit_status.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	if _, err = f.Write(protoStatus); err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readReports(w http.ResponseWriter, start string, end string) {
	var reports bit.TestResultsSet
	const layout = "2006-January-02 15:4:5"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	err = filepath.Walk("storage/test_reports/",
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && !reportTime.Before(startTime) && reportTime.Before(endTime) {
					protoReports, err := ioutil.ReadFile(path + "/tests_results.txt")
					if err != nil {
						rh.ApiResponseHandler(w, http.StatusInternalServerError, "Can't find report!", err)
					}
					var temp bit.TestResultsSet
					err = proto.Unmarshal(protoReports, &temp)
					if err != nil {
						rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
					reports.ResultsSet = append(reports.ResultsSet, temp.ResultsSet...)
				}
			}
			return nil
		})
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	response := make([]models.TestReport, len(reports.ResultsSet))
	for i, report := range reports.ResultsSet {
		response[i] = testResultToTestReport(*report)
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readConfigFailures(w http.ResponseWriter) {
	var configFailures []bit.Failure
	files, err := ioutil.ReadDir("storage/config/filtering_rules")
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	for _, f := range files {
		if f.Name() == ".gitignore" {
			continue
		}
		content, err := ioutil.ReadFile("storage/config/filtering_rules/" + f.Name())
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}

		decodedContent := bit.Failure{}
		err = proto.Unmarshal(content, &decodedContent)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}

		configFailures = append(configFailures, decodedContent)
	}
	err = json.NewEncoder(w).Encode(&configFailures)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readExtendedFailures(w http.ResponseWriter) {
	var foreverFailures []handler.ExtendedFailure
	files, err := ioutil.ReadDir("storage/config/perm_filtering_rules")
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	for _, f := range files {
		if f.Name() == ".gitignore" {
			continue
		}
		content, err := ioutil.ReadFile("storage/config/perm_filtering_rules/" + f.Name())
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		decodedContent := handler.ExtendedFailure{}
		err = proto.Unmarshal(content, &decodedContent)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		foreverFailures = append(foreverFailures, decodedContent)
	}
	err = json.NewEncoder(w).Encode(&foreverFailures)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readBitStatus(w http.ResponseWriter, start string, end string) {
	var statuses []bit.BitStatus
	const layout = "2006-January-02 15:4:5"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	err = filepath.Walk("storage/bit_status/",
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && !reportTime.Before(startTime) && reportTime.Before(endTime) {
					protoStatus, err := ioutil.ReadFile(path + "/bit_status.txt")
					if err != nil {
						rh.ApiResponseHandler(w, http.StatusInternalServerError, "Can't find status!", err)
					}
					temp := bit.BitStatus{}
					err = proto.Unmarshal(protoStatus, &temp)
					if err != nil {
						rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
					statuses = append(statuses, temp)
				}
			}
			return nil
		})
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	err = json.NewEncoder(w).Encode(&statuses)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readUserGroupMaskedTestIds(w http.ResponseWriter, userGroup string) {
	content, err := ioutil.ReadFile("storage/config/user_groups_masks/" + userGroup + ".txt")
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	decodedContent := bit.UserGroupsFiltering_FilteredFailures{}
	err = proto.Unmarshal(content, &decodedContent)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	err = json.NewEncoder(w).Encode(&decodedContent.MaskedTestIds)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readUserGroups(w http.ResponseWriter) {
	content, err := ioutil.ReadFile("storage/config/user_groups/user_groups.txt")
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	groups := strings.Split(string(content), "\n")
	err = json.NewEncoder(w).Encode(&groups)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteAgedData(w http.ResponseWriter, fileType string, timestamp string) {

	const layout = "2006-January-02 15:4:5"
	threshold, err := time.Parse(layout, timestamp)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	err = filepath.Walk("storage/"+fileType,
		func(path string, info os.FileInfo, err error) error {
			pathToTime := strings.Split(path, "\\")
			if len(pathToTime) >= 2 {
				timeToCmp := strings.ReplaceAll(pathToTime[2], " ", "-") + " " + strings.Join(pathToTime[3:], ":")
				reportTime, err := time.Parse(layout, timeToCmp)
				if err == nil && info.IsDir() && reportTime.Before(threshold) {
					err = os.RemoveAll(path)
					if err != nil {
						rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
					}
				}
			}
			return nil
		})
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	deleteEmptyDirectories(w, "storage/"+fileType, 0, fileType)
}

func deleteEmptyDirectories(w http.ResponseWriter, path string, depth int, fileType string) {
	if depth == 3 {
		return
	}
	err := filepath.Walk(path,
		func(childPath string, info fs.FileInfo, err error) error {
			deleteEmptyDirectories(w, childPath, depth+1, fileType)
			if (path != ("storage/" + fileType)) && IsEmpty(path) {
				err := os.RemoveAll(path)
				if err != nil {
					rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
				}
			}
			return nil
		})
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
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

func convertToKeyValuePair(arr []models.KeyValue) []*bit.KeyValuePair {
	copyArr := make([]*bit.KeyValuePair, len(arr))
	for k, v := range arr {
		pair := bit.KeyValuePair{Key: []byte(v.Key), Value: []byte(v.Value)}
		copyArr[k] = &pair
	}
	return copyArr
}

func convertToKeyValue(arr []*bit.KeyValuePair) []models.KeyValue {
	copyArr := make([]models.KeyValue, len(arr))
	for k, v := range arr {
		pair := models.KeyValue{Key: string(v.Key), Value: string(v.Value)}
		copyArr[k] = pair
	}
	return copyArr
}

func testReportToTestResult(tr models.TestReport) bit.TestResult {
	return bit.TestResult{
		TestId:         uint64(tr.TestId),
		Timestamp:      timestamppb.New(tr.Timestamp),
		TagSet:         convertToKeyValuePair(tr.TagSet),
		FieldSet:       convertToKeyValuePair(tr.FieldSet),
		ReportPriority: uint32(tr.ReportPriority),
	}
}

func testResultToTestReport(tr bit.TestResult) models.TestReport {
	return models.TestReport{
		TestId:         float64(tr.TestId),
		ReportPriority: float64(tr.ReportPriority),
		Timestamp:      tr.Timestamp.AsTime(),
		TagSet:         convertToKeyValue(tr.TagSet),
		FieldSet:       convertToKeyValue(tr.FieldSet),
	}
}
