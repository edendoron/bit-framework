package bithandler

import (
	"bytes"
	"encoding/json"
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	"github.com/edendoron/bit-framework/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var TestingFlag = false

// BitAnalyzer extract data from storage and creates bitStatus
type BitAnalyzer struct {
	ConfigFailures      []ExtendedFailure
	Reports             []models.TestReport
	LastEpochReportTime time.Time
	SavedFailures       []ExtendedFailure
	Status              bit.BitStatus
}

// ByTime implements sort.Interface for []TestReport based on the timestamp field.
type ByTime []models.TestReport

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }

// constant date layout used in storage
const layout = "2006-January-02 15:4:5"

// exported methods

func (a *BitAnalyzer) ReadFailuresFromStorage(keyValue string) {
	// read config failures
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		log.Printf("error create storage request. %v", err)
		return
	}

	params := req.URL.Query()
	params.Add(keyValue, "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	storageResponse, e := client.Do(req)
	if e != nil || storageResponse.StatusCode != http.StatusOK {
		log.Fatalln("error reading failures from storage")
	}
	defer storageResponse.Body.Close()

	switch keyValue {
	case "config_failures":
		var configFails []bit.Failure
		err = json.NewDecoder(storageResponse.Body).Decode(&configFails)
		a.ConfigFailures = FailuresSliceToExtendedFailuresSlice(configFails)
	case "forever_failures":
		err = json.NewDecoder(storageResponse.Body).Decode(&a.SavedFailures)
	}
	if err != nil {
		log.Printf("error reading %v from storage. %v", keyValue, err)
	}
}

func (a *BitAnalyzer) ReadReportsFromStorage() {
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		log.Printf("error create storage request %v", err)
		return
	}

	endTime := epoch.Format(layout)
	startTime := epoch.Add(-time.Second).Format(layout)
	params := req.URL.Query()
	params.Add("reports", "")
	params.Add("filter", "time")
	params.Add("start", startTime)
	params.Add("end", endTime)
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error execute storage request. %v", err)
		return
	}
	defer storageResponse.Body.Close()

	var storageReports []models.TestReport
	err = json.NewDecoder(storageResponse.Body).Decode(&storageReports)
	if err != nil {
		log.Printf("error decode storage response body. %v", err)
		return
	}

	// remove unnecessary reports and sort all reports by time
	a.UpdateReports(storageReports)
}

// Crosscheck examine reports and checks for failures violations
func (a *BitAnalyzer) Crosscheck(examinationTime time.Time) {

	// generate map of masked groups
	maskedUserGroups := make(map[string]int)

	// filter saved failures (remove NO_LATCH and NUM_OF_SECONDS if needed)
	a.filterSavedFailures()

	for idx, confFailure := range a.ConfigFailures {
		countFailed := a.checkExaminationRule(idx)

		if countFailed > 0 {
			// update masked groups map
			for _, userGroup := range confFailure.Failure.Dependencies.MasksOtherGroup {
				maskedUserGroups[userGroup] = 1
			}

			// insert failures to SavedFailures
			extFailure := confFailure
			extFailure.failureCount = countFailed
			extFailure.Time = examinationTime
			a.SavedFailures = append(a.SavedFailures, extFailure)

			//// post forever failures to storage in order to restore it later if needed
			if extFailure.Failure.ReportDuration.Indication == bit.FailureReportDuration_LATCH_FOREVER {
				writeForeverFailure(extFailure)
			}
		}
	}

	// update BitStatus (masked groups)
	a.updateBitStatus(maskedUserGroups)
}

// WriteBitStatus reports last interval bitStatus to storage
func (a *BitAnalyzer) WriteBitStatus() {

	if len(a.Status.Failures) == 0 {
		a.CleanBitStatus()
		return
	}

	jsonStatus, err := json.MarshalIndent(a.Status, "", " ")
	if err != nil {
		log.Printf("error marshal bit_status. %v", err)
		return
	}

	message := models.KeyValue{
		Key:   "bit_status",
		Value: string(jsonStatus),
	}

	jsonMessage, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		log.Printf("error marshal bit_status. %v", err)
		return
	}

	postBody := bytes.NewReader(jsonMessage)

	storageResponse, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error post bit_status to storage. %v", err)
	}
	storageResponse.Body.Close()

	a.CleanBitStatus()
}

// ResetSavedFailures handles user requests to remove LATCH_UNTIL_RESET duration failures from reported failures
func (a *BitAnalyzer) ResetSavedFailures() {
	n := 0
	for _, item := range a.SavedFailures {
		isResetIndication := item.Failure.ReportDuration.Indication == bit.FailureReportDuration_LATCH_UNTIL_RESET
		if !isResetIndication {
			// keep saved failure for next trigger check
			a.SavedFailures[n] = item
			n++
		}
	}
	a.SavedFailures = a.SavedFailures[:n]
}

// CleanBitStatus remove last interval bitStatus from analyzer after reporting it to storage
func (a *BitAnalyzer) CleanBitStatus() {
	a.Status = bit.BitStatus{}
}

// internal methods

// remove NO_LATCH and NUM_OF_SECONDS (if time interval has passed) indication failures from failure's list that needs to be reported in the next interval
func (a *BitAnalyzer) filterSavedFailures() {
	n := 0
	for _, item := range a.SavedFailures {
		indication := item.Failure.ReportDuration.Indication
		switch indication {
		case bit.FailureReportDuration_NO_LATCH:
		case bit.FailureReportDuration_NUM_OF_SECONDS:
			if uint32(time.Since(item.Time).Seconds()) < item.Failure.ReportDuration.IndicationSeconds {
				//keep saved failure for next trigger check
				a.SavedFailures[n] = item
				n++
			}
		default:
			// for until_reset and forever failures, keep saved failure for next trigger check
			a.SavedFailures[n] = item
			n++
		}
	}
	a.SavedFailures = a.SavedFailures[:n]
}

// checks for masked user groups in order to not report those failures in the bitStatus report.
func (a *BitAnalyzer) updateBitStatus(maskedUserGroups map[string]int) {
	for _, item := range a.SavedFailures {
		masked := checkMasked(maskedUserGroups, item)
		a.insertReportedFailureBitStatus(masked, item)
	}
}

// insert  un-masked failures to bitStatus report.
func (a *BitAnalyzer) insertReportedFailureBitStatus(masked bool, item ExtendedFailure) {
	if !masked {
		// insert failure to BitStatus
		reportedFailure := &bit.BitStatus_RportedFailure{
			FailureData: item.Failure.Description,
			Timestamp:   timestamppb.New(item.Time),
			Count:       item.failureCount,
		}
		a.Status.Failures = append(a.Status.Failures, reportedFailure)
	}
}

/*
checks failure violation.
@param i is the index of the tested failure in analyzer config failures array.
@returns the count of failure violation found in current time frame
*/
func (a *BitAnalyzer) checkExaminationRule(i int) uint64 {
	examinationRule := a.ConfigFailures[i].Failure.ExaminationRule
	timeCriteria := examinationRule.FailureCriteria.TimeCriteria
	if len(a.Reports) == 0 {
		return 0
	}
	// count failed tests
	switch timeCriteria.WindowType {
	case bit.FailureExaminationRule_FailureCriteria_FailureTimeCriteria_NO_WINDOW:
		return a.checkNoWindow(i)
	case bit.FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING:
		return a.checkSlidingWindow(i)
	default:
		return 0
	}
}

/*
iterate reports in order to check failure violation for SLIDING_WINDOW type of failures
@param i is the index of the tested failure in analyzer config failures array.
@returns the count of failure violation found in current time frame
*/
func (a *BitAnalyzer) checkSlidingWindow(i int) uint64 {

	extFailure := &a.ConfigFailures[i]
	examinationRule := extFailure.Failure.ExaminationRule
	timeCriteria := examinationRule.FailureCriteria.TimeCriteria
	var countExaminationRuleViolation uint64 = 0
	countTimeCriteriaViolation := extFailure.reportsCount

	begin := len(a.Reports)
	for idx := range a.Reports {
		if a.Reports[idx].TestId == extFailure.startReportId {
			begin = idx
			break
		}
	}
	end := len(a.Reports) + 1
	for idx := range a.Reports {
		if a.Reports[idx].TestId == extFailure.endReportId {
			end = idx + 1
			break
		}
	}

	if end == len(a.Reports) {
		return 0
	}

	if begin == len(a.Reports) || end == len(a.Reports)+1 {
		begin = 0
		end = 0
		countTimeCriteriaViolation = 0
	}
	for end < len(a.Reports) {
		startWindowTest := a.Reports[begin]
		endWindowTest := a.Reports[end]

		timeDiff := endWindowTest.Timestamp.Sub(startWindowTest.Timestamp)
		if timeDiff.Seconds() <= float64(timeCriteria.WindowSize) {
			if failedValueCriteria(&endWindowTest, examinationRule) {
				countTimeCriteriaViolation++
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

	end--
	//set new fields in order to continue from the same point in the next epoch
	extFailure.startReportId = a.Reports[begin].TestId
	extFailure.endReportId = a.Reports[end].TestId
	extFailure.reportsCount = countTimeCriteriaViolation

	// set new LastEpochReportTime
	if a.Reports[begin].Timestamp.Before(a.LastEpochReportTime) {
		a.LastEpochReportTime = a.Reports[begin].Timestamp
	}

	return countExaminationRuleViolation
}

/*
iterate reports in order to check failure violation for NO_WINDOW type of failures
@param i is the index of the tested failure in analyzer config failures array.
@returns the count of failure violation found in current time frame
*/
func (a *BitAnalyzer) checkNoWindow(i int) uint64 {
	extFailure := &a.ConfigFailures[i]
	examinationRule := extFailure.Failure.ExaminationRule

	var countExaminationRuleViolation uint64 = 0
	begin := len(a.Reports) + 1
	for idx := range a.Reports {
		if a.Reports[idx].TestId == extFailure.endReportId {
			begin = idx + 1
			break
		}
	}
	if begin == len(a.Reports)+1 {
		begin = 0
	}
	for begin < len(a.Reports) {
		if failedValueCriteria(&a.Reports[begin], examinationRule) {
			countExaminationRuleViolation++
		}
		begin++
	}

	begin--
	// set start and end report Id's to last report in order to avoid duplicates on the next epoch, note that reportsCount is always 0 for noWindow failures
	extFailure.startReportId = a.Reports[len(a.Reports)-1].TestId
	extFailure.endReportId = extFailure.startReportId

	// set new LastEpochReportTime
	if a.Reports[begin].Timestamp.Before(a.LastEpochReportTime) {
		a.LastEpochReportTime = a.Reports[begin].Timestamp
	}

	return countExaminationRuleViolation
}

/*
UpdateReports removes previous time frame reports from analyzer reports, update it with current time frame reports and sort them by time.
note: some of previous time frame reports may not be removed in order to keep examine failures from the last point stopped,
and find failures that may occur between time frames (relevant for SLIDING_WINDOW failures).
*/
func (a *BitAnalyzer) UpdateReports(reports []models.TestReport) {
	n := sort.Search(len(a.Reports), func(i int) bool { return !a.Reports[i].Timestamp.Before(a.LastEpochReportTime) })
	a.Reports = append(a.Reports[n:], reports...)
	sort.Stable(ByTime(a.Reports))
	if len(a.Reports) > 0 {
		a.LastEpochReportTime = a.Reports[len(a.Reports)-1].Timestamp
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

// return true iff value criteria is failed due to field or tag matching problem.
func failedValueCriteria(test *models.TestReport, examinationRule *bit.FailureExaminationRule) bool {

	fieldValue := checkField(test, examinationRule)
	if fieldValue == "" || !checkTag(test, examinationRule) {
		return false
	}

	return checkFailedValue(fieldValue, examinationRule.FailureCriteria.ValueCriteria)
}

// check field existing
func checkField(test *models.TestReport, examinationRule *bit.FailureExaminationRule) string {
	for _, field := range test.FieldSet {
		if field.Key == examinationRule.MatchingField {
			return field.Value
		}
	}
	return ""
}

// check tag existing
func checkTag(test *models.TestReport, examinationRule *bit.FailureExaminationRule) bool {
	// fixes problem of base64 encoding of protobuf
	key, err := json.MarshalIndent(examinationRule.MatchingTag.Key, "", " ")
	if err != nil {
		log.Printf("error mrashalling tag set of examination rule error: %v", err)
	}
	value, err := json.MarshalIndent(examinationRule.MatchingTag.Value, "", " ")
	if err != nil {
		log.Printf("error mrashalling tag set of examination rule error: %v", err)
	}
	if TestingFlag {
		key = examinationRule.MatchingTag.Key
		value = examinationRule.MatchingTag.Value
	}
	for _, tag := range test.TagSet {
		if tag.Key == strings.ReplaceAll(string(key), "\"", "") && tag.Value == strings.ReplaceAll(string(value), "\"", "") {
			return true
		}
	}
	return false
}

// check failure value criteria, return true if value of test violates rule, false otherwise
func checkFailedValue(value string, criteria *bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria) bool {
	min, max := calculateThreshold(criteria.Minimum, criteria.Miximum, criteria)
	if min >= max {
		return false
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("error convert string value %v to number", value)
		return false
	}
	valueWithin := floatValue >= min && floatValue <= max
	switch criteria.ThresholdMode {
	case bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN:
		return valueWithin
	case bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF:
		return !valueWithin
	default:
		return valueWithin
	}
}

// calculate threshold after allowed deviation (by value or percent) defined in failure configuration. return (min, max) threshold
func calculateThreshold(minimum float64, maximum float64, criteria *bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria) (float64, float64) {
	deviation := criteria.Exceeding.Value
	if criteria.Exceeding.Type == bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT {
		deviation *= 0.01 * (maximum - minimum)
	}
	if criteria.ThresholdMode == bit.FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF {
		return minimum - deviation, maximum + deviation
	} else {
		return minimum + deviation, maximum - deviation
	}
}

// when analyzer finds violation of a failure defined as FOREVER_FAILURE, writeForeverFailure post this failure to storage in order to restore it even if the service is restarted
func writeForeverFailure(failure ExtendedFailure) {
	if TestingFlag {
		return
	}
	jsonForeverFailure, err := json.MarshalIndent(failure, "", " ")
	if err != nil {
		log.Printf("error marshal forever_failure. %v", err)
		return
	}

	message := models.KeyValue{
		Key:   "forever_failure",
		Value: string(jsonForeverFailure),
	}

	jsonMessage, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		log.Printf("error marshal forever_failure message. %v", err)
		return
	}
	postBody := bytes.NewReader(jsonMessage)

	storageResponse, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBody)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		log.Printf("error post forever_failure to storage. %v", err)
	}
	storageResponse.Body.Close()
}
