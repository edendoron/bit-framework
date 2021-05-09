package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
)

type Analyzer struct {
	Failures []Failure
	Reports  []TestResult
	Status   BitStatus
}

func (a *Analyzer) ReadFromStorage(key string) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		// TODO: handle error
		return
	}
	//defer req.Body.Close()

	params := req.URL.Query()
	params.Add("key", key)

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		// TODO: handle error
		return
	}
	defer storageResponse.Body.Close()

	if key == "Failure" {
		err = json.NewDecoder(storageResponse.Body).Decode(&a.Failures)
	} else if key == "Report" {
		err = json.NewDecoder(storageResponse.Body).Decode(&a.Reports)
	}

	if err != nil {
		// TODO: handle error
		return
	}
}

func (a *Analyzer) Crosscheck() {
	for _, failure := range a.Failures {
		countFailed, timestamp := a.checkExaminationRule(failure)
		if countFailed > 0 {
			reportedFailure := &BitStatus_RportedFailure{
				FailureData: failure.Description,
				Timestamp:   timestamp,
				Count:       countFailed,
			}
			a.Status.Failures = append(a.Status.Failures, reportedFailure)
		}
		// TODO: a.Status.UserGroup
	}
}

func (a *Analyzer) WriteBitStatus() {

	//TODO: handle error
	jsonStatus, _ := json.MarshalIndent(a.Status, "", " ")

	message := KeyValuePair{
		Key:   []byte("BitStatus"),
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
}

func (a *Analyzer) checkExaminationRule(failure Failure) (uint64, *timestamppb.Timestamp) {
	var countExaminationRuleViolation uint64 = 0
	time := timestamppb.Now()
	if len(a.Reports) == 0 {
		return countExaminationRuleViolation, time
	}
	examinationRule := failure.ExaminationRule
	failureCriteria := examinationRule.FailureCriteria
	timeCriteria := failureCriteria.TimeCriteria
	// count failed tests
	var countTimeCriteriaViolation uint32 = 0
	switch timeCriteria.WindowType {
	case FailureExaminationRule_FailureCriteria_FailureTimeCriteria_NO_WINDOW:
		for _, test := range a.Reports {
			if failedValueCriteria(&test, examinationRule) {
				countTimeCriteriaViolation++
				time = test.Timestamp
			}
			if countTimeCriteriaViolation > timeCriteria.FailuresCCount {
				countExaminationRuleViolation = 1
			}
		}
	case FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING:
		// we assume that reports from storage are sorted by time
		begin, end := 0, 0
		for end < len(a.Reports) {
			startWindowTest := a.Reports[begin]
			endWindowTest := a.Reports[end]
			timeDiff := endWindowTest.Timestamp.Seconds - startWindowTest.Timestamp.Seconds
			if timeDiff <= int64(timeCriteria.WindowSize) {
				if failedValueCriteria(&endWindowTest, examinationRule) {
					countTimeCriteriaViolation++
					time = endWindowTest.Timestamp
				}
				end++
				if end == len(a.Reports) && countTimeCriteriaViolation > timeCriteria.FailuresCCount {
					countExaminationRuleViolation++
				}
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
	}
	return countExaminationRuleViolation, time
}

func failedValueCriteria(test *TestResult, examinationRule *FailureExaminationRule) bool {

	fieldValue := checkField(test, examinationRule)
	if fieldValue == "" || !checkTag(test, examinationRule) {
		return false
	}

	return checkValue(fieldValue, examinationRule.FailureCriteria.ValueCriteria)
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
func checkValue(value string, criteria *FailureExaminationRule_FailureCriteria_FailureValueCriteria) bool {
	min, max := calculateThreshold(criteria.Minimum, criteria.Miximum, criteria.Exceeding)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		//TODO: handle error
		return false
	}
	valueWithin := floatValue >= min && floatValue <= max
	switch criteria.ThresholdMode {
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN:
		return valueWithin
	case FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF:
		return !valueWithin
	default:
		return valueWithin
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
