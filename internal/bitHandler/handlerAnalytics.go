package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BitAnalyzer struct {
	ConfigFailures []Failure
	Reports        []TestReport
	SavedFailures  []ExtendedFailure
	Status         BitStatus
}

type ExtendedFailure struct {
	Failure Failure
	time    time.Time
	count   uint64
}

func (e *ExtendedFailure) Reset() { *e = ExtendedFailure{} }

func (e *ExtendedFailure) String() string { return proto.CompactTextString(e) }

func (e *ExtendedFailure) ProtoMessage() {}

// exported methods

func (a *BitAnalyzer) ReadFailuresFromStorage(keyValue string) {
	// read config failures
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		log.Printf("error create storage request")
		return
	}
	//TODO: defer req.Body.Close()

	params := req.URL.Query()
	params.Add(keyValue, "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	storageResponse, e := client.Do(req)
	if e != nil || storageResponse.StatusCode != http.StatusOK {
		log.Fatalln("error reading failures from storage")
	}

	switch keyValue {
	case "config_failures":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.ConfigFailures)
	case "forever_failures":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.SavedFailures)
	}
	if err != nil {
		log.Println("error reading " + keyValue + " from storage")
	}

	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}
}

func (a *BitAnalyzer) ReadReportsFromStorage(d time.Duration) {
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		log.Printf("error create storage request")
		return
	}
	//defer req.Body.Close()

	const layout = "2006-January-02 15:4:5"
	//TODO: update time frame according to d
	endTime, _ := time.Parse(layout, "2021-April-15 12:00:00")
	startTime, _ := time.Parse(layout, "2021-April-15 11:00:00")
	params := req.URL.Query()
	params.Add("reports", "")
	params.Add("filter", "time")
	params.Add("start", startTime.Format(layout))
	params.Add("end", endTime.Format(layout))
	req.URL.RawQuery = params.Encode()

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
			for _, userGroup := range failure.Dependencies.MasksOtherGroup {
				maskedUserGroups[userGroup] = 1
			}

			// insert failures to SavedFailures
			timedFailure := ExtendedFailure{
				Failure: failure,
				time:    timestamp,
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

	storageResponse, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBody)
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
		isResetIndication := item.Failure.ReportDuration.Indication == FailureReportDuration_LATCH_UNTIL_RESET
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
		indication := item.Failure.ReportDuration.Indication
		switch indication {
		case FailureReportDuration_NO_LATCH:
			a.insertReportedFailureBitStatus(masked, item)
		case FailureReportDuration_NUM_OF_SECONDS:
			if uint32(time.Since(item.time)) < item.Failure.ReportDuration.IndicationSeconds {
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

func (a *BitAnalyzer) insertReportedFailureBitStatus(masked bool, item ExtendedFailure) {
	if !masked {
		// insert failure to BitStatus
		reportedFailure := &BitStatus_RportedFailure{
			FailureData: item.Failure.Description,
			Timestamp:   timestamppb.New(item.time),
			Count:       item.count,
		}
		a.Status.Failures = append(a.Status.Failures, reportedFailure)
	}
}

// return true iff all of the BelongsToGroup are masked by other failures
func checkMasked(maskedUserGroups map[string]int, item ExtendedFailure) bool {
	countBelongGroupMask := 0
	for _, group := range item.Failure.Dependencies.BelongsToGroup {
		if maskedUserGroups[group] == 1 {
			countBelongGroupMask++
		}
	}
	return countBelongGroupMask == len(item.Failure.Dependencies.BelongsToGroup)
}

func (a *BitAnalyzer) checkExaminationRule(failure Failure) (uint64, time.Time) {
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
		return 0, time.Now()
	}
}

func (a *BitAnalyzer) checkSlidingWindow(timeCriteria *FailureExaminationRule_FailureCriteria_FailureTimeCriteria, examinationRule *FailureExaminationRule) (uint64, time.Time) {
	var countExaminationRuleViolation uint64 = 0
	timestamp := time.Now()
	var countTimeCriteriaViolation uint32 = 0
	begin, end := 0, 0
	for end < len(a.Reports) {
		startWindowTest := a.Reports[begin]
		endWindowTest := a.Reports[end]

		timeDiff := endWindowTest.Timestamp.Sub(startWindowTest.Timestamp)
		if timeDiff.Seconds() <= float64(timeCriteria.WindowSize) {
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

func (a *BitAnalyzer) checkNoWindow(examinationRule *FailureExaminationRule, timeCriteria *FailureExaminationRule_FailureCriteria_FailureTimeCriteria) (uint64, time.Time) {
	var countExaminationRuleViolation uint64 = 0
	timestamp := time.Now()
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

func failedValueCriteria(test *TestReport, examinationRule *FailureExaminationRule) bool {

	fieldValue := checkField(test, examinationRule)
	if fieldValue == "" || !checkTag(test, examinationRule) {
		return false
	}

	return checkFailedValue(fieldValue, examinationRule.FailureCriteria.ValueCriteria)
}

// check field existing
func checkField(test *TestReport, examinationRule *FailureExaminationRule) string {
	for _, field := range test.FieldSet {
		if string(field.Key) == examinationRule.MatchingField {
			return string(field.Value)
		}
	}
	return ""
}

// check tag existing
func checkTag(test *TestReport, examinationRule *FailureExaminationRule) bool {
	//for _, tag := range test.TagSet {
	//	if tag.Key == string(examinationRule.MatchingTag.Key) && tag.Value == string(examinationRule.MatchingTag.Value) {
	//		return true
	//	}
	//}
	//return false
	return true
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

	storageResponse, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error post forever_failure to storage")
		return
	}
	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
	}
}
