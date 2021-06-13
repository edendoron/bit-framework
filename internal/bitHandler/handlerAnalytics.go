package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BitAnalyzer struct {
	ConfigFailures []Failure
	Reports        []TestResult
	SavedFailures  []extendedFailure
	Status         BitStatus
}

type extendedFailure struct {
	failure Failure
	time    time.Time
	count   uint64
}

// exported methods

//func (a *BitAnalyzer) ReadFailureFromLocalConfigFile() {
//	failure := Failure{}
//	content, err := ioutil.ReadFile("./configs/config_failures/voltage_failure.json")
//	if err != nil {
//		//TODO: handle error
//	}
//	err = json.Unmarshal(content, &failure)
//	if err != nil {
//		//TODO: handle error
//	}
//	failure.ExaminationRule.MatchingTag.Key = []byte("zone")
//	failure.ExaminationRule.MatchingTag.Value = []byte("north")
//
//	a.ConfigFailures = append(a.ConfigFailures, failure)
//}
//
//func (a *BitAnalyzer) ReadReportsFromLocalConfigFile() {
//	reportBody := TestReport{}
//	content, err := ioutil.ReadFile("./storage/reports/reports.json")
//	if err != nil {
//		//TODO: handle error
//	}
//	err = json.Unmarshal(content, &reportBody)
//	if err != nil {
//		//TODO: handle error
//	}
//	tagset1 := KeyValuePair{
//		Key:   []byte("zone"),
//		Value: []byte("north"),
//	}
//	tagset2 := KeyValuePair{
//		Key:   []byte("hostname"),
//		Value: []byte("server02"),
//	}
//	fieldset1 := KeyValuePair{
//		Key:   []byte("TemperatureCelsius"),
//		Value: []byte("-40.8"),
//	}
//	fieldset2 := KeyValuePair{
//		Key:   []byte("volts"),
//		Value: []byte("7.8"),
//	}
//	testResult := TestResult{
//		TestId:         uint64(reportBody.TestId),
//		Timestamp:      timestamppb.New(reportBody.Timestamp),
//		TagSet:         []*KeyValuePair{&tagset1, &tagset2},
//		FieldSet:       []*KeyValuePair{&fieldset1, &fieldset2},
//		ReportPriority: uint32(reportBody.ReportPriority),
//	}
//	a.Reports = append(a.Reports, testResult)
//
//}

func (a *BitAnalyzer) ReadFailuresFromStorage(keyValue string) {
	// read config failures
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		log.Printf("error create storage request")
		return
	}
	//TODO: defer req.Body.Close()

	params := req.URL.Query()
	params.Add(keyValue, "")

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}

	switch keyValue {
	case "config_failure":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.ConfigFailures)
	case "forever_failure":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.SavedFailures)
	}
	if err != nil {
		log.Println("error reading failures from storage")
		return
	}

	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}
}

func (a *BitAnalyzer) ReadReportsFromStorage(d time.Duration) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		log.Printf("error create storage request")
		return
	}
	//defer req.Body.Close()

	endTime := time.Now()
	startTime := endTime.Add(-d)
	params := req.URL.Query()
	params.Add("report", "")
	params.Add("filter", "time")
	params.Add("start", startTime.String())
	params.Add("end", endTime.String())

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}

	err = json.NewDecoder(storageResponse.Body).Decode(&a.Reports)
	if err != nil {
		log.Printf("error decode storage response body")
		return
	}
	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}

}

func (a *BitAnalyzer) Crosscheck() {

	// generate map of masked groups
	maskedUserGroups := make(map[string]int)

	for _, failure := range a.ConfigFailures {
		countFailed, timestamp := a.checkExaminationRule(failure)

		if countFailed > 0 {
			// update masked groups map
			for _, userGroup := range failure.Dependencies.BelongsToGroup {
				maskedUserGroups[userGroup] = 1
			}

			// insert failures to SavedFailures
			timedFailure := extendedFailure{
				failure: failure,
				time:    timestamp.AsTime(),
				count:   countFailed,
			}
			a.SavedFailures = append(a.SavedFailures, timedFailure)

			// post forever failures to storage in order to restore it later if needed
			if failure.ReportDuration.Indication == FailureReportDuration_LATCH_FOREVER {
				writeForeverFailure(failure)
			}
		}
	}

	// filter saved failures and update BitStatus (timed failures and masked groups)
	a.updateBitStatus(maskedUserGroups)
}

func (a *BitAnalyzer) WriteBitStatus() {

	jsonStatus, err := json.MarshalIndent(a.Status, "", " ")
	if err != nil {
		log.Printf("error marshal bit_status")
		return
	}

	message := KeyValue{
		Key:   "bit_status",
		Value: string(jsonStatus),
	}

	jsonMessage, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		log.Printf("error marshal bit_status")
		return
	}

	postBody := bytes.NewReader(jsonMessage)

	storageResponse, err := http.Post(storageDataWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error post bit_status to storage")
		return
	}
	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}

	a.cleanBitStatus()
}

func (a *BitAnalyzer) FilterSavedFailures() {
	n := 0
	for _, item := range a.SavedFailures {
		isResetIndication := item.failure.ReportDuration.Indication == FailureReportDuration_LATCH_UNTIL_RESET
		if !isResetIndication {
			// keep saved failure for next trigger check
			a.SavedFailures[n] = item
			n++
		}
	}
	a.SavedFailures = a.SavedFailures[:n]
}

// internal methods

func (a *BitAnalyzer) cleanBitStatus() {
	a.Status = BitStatus{}
}

func (a *BitAnalyzer) updateBitStatus(maskedUserGroups map[string]int) {
	n := 0
	for _, item := range a.SavedFailures {
		masked := checkMasked(maskedUserGroups, item)
		indication := item.failure.ReportDuration.Indication
		switch indication {
		case FailureReportDuration_NO_LATCH:
			a.insertReportedFailureBitStatus(masked, item)
		case FailureReportDuration_NUM_OF_SECONDS:
			if uint32(time.Since(item.time)) < item.failure.ReportDuration.IndicationSeconds {
				//keep saved failure for next trigger check
				a.SavedFailures[n] = item
				n++

				a.insertReportedFailureBitStatus(masked, item)
			}
		default:
			// for until_reset and forever failures
			//keep saved failure for next trigger check
			a.SavedFailures[n] = item
			n++

			a.insertReportedFailureBitStatus(masked, item)
		}
	}
	a.SavedFailures = a.SavedFailures[:n]
}

func (a *BitAnalyzer) insertReportedFailureBitStatus(masked bool, item extendedFailure) {
	if !masked {
		// insert failure to BitStatus
		reportedFailure := &BitStatus_RportedFailure{
			FailureData: item.failure.Description,
			Timestamp:   timestamppb.New(item.time),
			Count:       item.count,
		}
		a.Status.Failures = append(a.Status.Failures, reportedFailure)
	}
}

func checkMasked(maskedUserGroups map[string]int, item extendedFailure) bool {
	countBelongGroupMask := 0
	for _, group := range item.failure.Dependencies.BelongsToGroup {
		if maskedUserGroups[group] == 1 {
			countBelongGroupMask++
		}
	}
	return countBelongGroupMask == len(item.failure.Dependencies.BelongsToGroup)
}

func (a *BitAnalyzer) checkExaminationRule(failure Failure) (uint64, *timestamppb.Timestamp) {
	examinationRule := failure.ExaminationRule
	timeCriteria := examinationRule.FailureCriteria.TimeCriteria
	// count failed tests
	switch timeCriteria.WindowType {
	case FailureExaminationRule_FailureCriteria_FailureTimeCriteria_NO_WINDOW:
		return a.checkNoWindow(examinationRule, timeCriteria)
	case FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING:
		// TODO: make sure that reports from storage are sorted by timestamp
		return a.checkSlidingWindow(timeCriteria, examinationRule)
	default:
		return 0, timestamppb.Now()
	}
}

func (a *BitAnalyzer) checkSlidingWindow(timeCriteria *FailureExaminationRule_FailureCriteria_FailureTimeCriteria, examinationRule *FailureExaminationRule) (uint64, *timestamppb.Timestamp) {
	var countExaminationRuleViolation uint64 = 0
	timestamp := timestamppb.Now()
	var countTimeCriteriaViolation uint32 = 0
	begin, end := 0, 0
	for end < len(a.Reports) {
		startWindowTest := a.Reports[begin]
		endWindowTest := a.Reports[end]

		timeDiff := endWindowTest.Timestamp.Seconds - startWindowTest.Timestamp.Seconds
		if timeDiff <= int64(timeCriteria.WindowSize) {
			if failedValueCriteria(&endWindowTest, examinationRule) {
				countTimeCriteriaViolation++
				timestamp = endWindowTest.Timestamp
			}
			end++
		} else {
			if countTimeCriteriaViolation > timeCriteria.FailuresCCount {
				countExaminationRuleViolation++
			}
			if failedValueCriteria(&startWindowTest, examinationRule) {
				countTimeCriteriaViolation--
			}
			begin++
		}
	}
	// check for end of reports (last report in time frame)
	if countTimeCriteriaViolation > timeCriteria.FailuresCCount {
		countExaminationRuleViolation++
	}
	return countExaminationRuleViolation, timestamp
}

func (a *BitAnalyzer) checkNoWindow(examinationRule *FailureExaminationRule, timeCriteria *FailureExaminationRule_FailureCriteria_FailureTimeCriteria) (uint64, *timestamppb.Timestamp) {
	var countExaminationRuleViolation uint64 = 0
	timestamp := timestamppb.Now()
	var countTimeCriteriaViolation uint32 = 0
	for _, test := range a.Reports {
		if failedValueCriteria(&test, examinationRule) {
			countTimeCriteriaViolation++
			timestamp = test.Timestamp
		}
		if countTimeCriteriaViolation > timeCriteria.FailuresCCount {
			countExaminationRuleViolation = 1
		}
	}
	return countExaminationRuleViolation, timestamp
}

func failedValueCriteria(test *TestResult, examinationRule *FailureExaminationRule) bool {

	fieldValue := checkField(test, examinationRule)
	if fieldValue == "" || !checkTag(test, examinationRule) {
		return false
	}

	return checkFailedValue(fieldValue, examinationRule.FailureCriteria.ValueCriteria)
}

// check field existing
func checkField(test *TestResult, examinationRule *FailureExaminationRule) string {
	for _, field := range test.FieldSet {
		if string(field.Key) == examinationRule.MatchingField {
			return string(field.Value)
		}
	}
	return ""
}

// check tag existing
func checkTag(test *TestResult, examinationRule *FailureExaminationRule) bool {
	for _, tag := range test.TagSet {
		if string(tag.Key) == string(examinationRule.MatchingTag.Key) && string(tag.Value) == string(examinationRule.MatchingTag.Value) {
			return true
		}
	}
	return false
}

// check failure value criteria
func checkFailedValue(value string, criteria *FailureExaminationRule_FailureCriteria_FailureValueCriteria) bool {
	min, max := calculateThreshold(criteria.Minimum, criteria.Miximum, criteria.Exceeding)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		//TODO: handle error
		return false
	}
	valueWithin := floatValue >= min && floatValue <= max
	switch criteria.ThresholdMode {
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN:
		return !valueWithin
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF:
		return valueWithin
	default:
		return !valueWithin
	}
}

func calculateThreshold(minimum float64, maximum float64, exceeding *FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding) (float64, float64) {
	switch exceeding.Type {
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_VALUE:
		return minimum - exceeding.Value, maximum + exceeding.Value
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT:
		percent := exceeding.Value * 0.01
		return minimum * (1 - percent), maximum * (1 + percent)
	default:
		return minimum, maximum
	}
}

func writeForeverFailure(failure Failure) {
	jsonForeverFailure, err := json.MarshalIndent(failure, "", " ")
	if err != nil {
		log.Printf("error marshal forever_failure")
		return
	}

	message := KeyValue{
		Key:   "forever_failure",
		Value: string(jsonForeverFailure),
	}

	jsonMessage, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		log.Printf("error marshal forever_failure")
		return
	}
	postBody := bytes.NewReader(jsonMessage)

	storageResponse, err := http.Post(storageDataWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error post forever_failure to storage")
		return
	}
	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}
}
