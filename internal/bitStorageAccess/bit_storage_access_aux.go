package bitStorageAccess

import (
	. "../../configs/rafael.com/bina/bit"
	. "../apiResponseHandlers"
	. "../models"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"os"
)


func postReports(w http.ResponseWriter, reports *string) {
	request := TestReport{}
	if err := json.Unmarshal([]byte(*reports), &request); err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	requestToProto := testReportToTestResult(request)
	protoReports, err := proto.Marshal(&requestToProto)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	path := "../storage/" + fmt.Sprint(request.Timestamp.Date()) + "/" + fmt.Sprint(request.Timestamp.Hour()) +
		"/" + fmt.Sprint(request.Timestamp.Minute()) + "/" + fmt.Sprint(request.Timestamp.Second())
	if _, err = os.Stat(path + "/tests_results.txt"); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		}
	}
	f, err := os.OpenFile(path + "/tests_results.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil{
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	if _, err = f.Write(protoReports); err != nil{
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
