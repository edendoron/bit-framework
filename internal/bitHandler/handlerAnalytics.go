package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
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
//	content, err := ioutil.ReadFile("./configs/config_failures/voltage_week_failure.json")
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
		// TODO: handle error
		return
	}
	//defer req.Body.Close()

	params := req.URL.Query()
	params.Add("key", keyValue)

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}
	defer storageResponse.Body.Close()

	switch keyValue {
	case "config_failure":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.ConfigFailures)
	case "forever_failure":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.SavedFailures)
	}
	if err != nil {
		// TODO: handle error
		return
	}
}

func (a *BitAnalyzer) ReadReportsFromStorage(d time.Duration) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		// TODO: handle error
		return
	}
	//defer req.Body.Close()

	params := req.URL.Query()
	params.Add("key", "report")
	params.Add("duration", d.String())

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}
	defer storageResponse.Body.Close()

	err = json.NewDecoder(storageResponse.Body).Decode(&a.Reports)
	if err != nil {
		// TODO: handle error
		return
	}
}

func (a *BitAnalyzer) Crosscheck() {
	for _, failure := range a.ConfigFailures {
		countFailed, timestamp := a.checkExaminationRule(failure)

		// insert old failures to BitStatus
		a.filterAndUpdateFailures()

		if countFailed > 0 {
			// insert failure to BitStatus
			reportedFailure := &BitStatus_RportedFailure{
				FailureData: failure.Description,
				//TODO: check if need to change to timestamppb.Now()
				Timestamp: timestamp,
				Count:     countFailed,
			}
			a.Status.Failures = append(a.Status.Failures, reportedFailure)

			// insert timed failures (of any kind) to SavedFailures
			if failure.ReportDuration.Indication != FailureReportDuration_NO_LATCH {
				timedFailure := extendedFailure{
					failure: failure,
					time:    timestamp.AsTime(),
					count:   countFailed,
				}
				a.SavedFailures = append(a.SavedFailures, timedFailure)
			}
		}
		// TODO: a.Status.UserGroup
	}
}

func (a *BitAnalyzer) WriteBitStatus() {

	//TODO: handle error
	jsonStatus, _ := json.MarshalIndent(a.Status, "", " ")

	message := KeyValuePair{
		Key:   []byte("bit_status"),
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

	a.cleanBitStatus()
}

// internal methods

func (a *BitAnalyzer) cleanBitStatus() {
	a.Status = BitStatus{}
}

func (a *BitAnalyzer) filterAndUpdateFailures() {
	n := 0
	for _, item := range a.SavedFailures {
		isSecondIndication := item.failure.ReportDuration.Indication == FailureReportDuration_NUM_OF_SECONDS
		if !isSecondIndication || uint32(time.Since(item.time)) < item.failure.ReportDuration.IndicationSeconds {

			// insert failure to BitStatus
			reportedFailure := &BitStatus_RportedFailure{
				FailureData: item.failure.Description,
				Timestamp:   timestamppb.New(item.time),
				Count:       item.count,
			}
			a.Status.Failures = append(a.Status.Failures, reportedFailure)

			// keep saved failure for next trigger check
			a.SavedFailures[n] = item
			n++
		}
	}
	a.SavedFailures = a.SavedFailures[:n]
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
