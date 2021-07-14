package bitHandler

import (
	. ".."
	. "../../../configs/rafael.com/bina/bit"
	. "../../models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
	"time"
)

// test for filter failures
func TestResetSavedFailures(t *testing.T) {

	fail0 := ExtendedFailure{Failure: Failure{ReportDuration: &FailureReportDuration{Indication: 3}}}
	fail1 := ExtendedFailure{Failure: Failure{ReportDuration: &FailureReportDuration{Indication: 1}}}
	fail2 := ExtendedFailure{Failure: Failure{ReportDuration: &FailureReportDuration{Indication: 0}}}
	fail3 := ExtendedFailure{Failure: Failure{ReportDuration: &FailureReportDuration{Indication: 2}}}
	fail4 := ExtendedFailure{Failure: Failure{ReportDuration: &FailureReportDuration{Indication: 1}}}

	fails := []ExtendedFailure{fail0, fail1, fail2, fail3, fail4}
	expectedResult := []ExtendedFailure{fail0, fail2, fail3}

	var analyzer BitAnalyzer
	analyzer.SavedFailures = fails

	analyzer.ResetSavedFailures()

	if !reflect.DeepEqual(analyzer.SavedFailures, expectedResult) {
		t.Errorf("Saved Failures did not filter properly")
	}

}

// test for update reports
func TestUpdateReports(t *testing.T) {

	var analyzer BitAnalyzer
	analyzer.Reports = []TestReport{report6, report7, report8, report9, report10}
	analyzer.LastEpochReportTime = report9.Timestamp

	analyzer.UpdateReports([]TestReport{report12, report11})

	expectedResult := []TestReport{report8, report9, report10, report11, report12}

	if !reflect.DeepEqual(analyzer.Reports, expectedResult) {
		t.Errorf("UpdateReports did not filter properly")
	}

}

// test for Crosscheck function, examination rule, reliability of analysis
func TestExaminationRule(t *testing.T) {

	var a BitAnalyzer

	TestingFlag = true

	testTime := time.Now()
	expectedResult := BitStatus{}

	// test empty bitStatus - no failures found
	a.ConfigFailures = []ExtendedFailure{failure1}
	a.Reports = []TestReport{report0}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test empty bitStatus, expected no failures, got failures")
	}
	clearBitStatus(&a)

	// test empty bitStatus - no reports in time frame
	a.ConfigFailures = []ExtendedFailure{failure1}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test empty bitStatus, expected no failures, got failures")
	}
	clearBitStatus(&a)

	// test failure Field + Tag match - one failure should be reported + WITHIN range threshold + SLIDING - one report violation
	a.ConfigFailures = []ExtendedFailure{failure0}
	a.Reports = []TestReport{report0}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure0.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       1,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure Field + Tag match, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
		if !reflect.DeepEqual(a.Status.Failures, expectedResult.Failures) {
			t.Errorf("test failure Field + Tag match, found failures diff")
		}
	}
	clearBitStatus(&a)

	// test failure only Field match - no failures should be reported
	a.ConfigFailures = []ExtendedFailure{failure0}
	a.Reports = []TestReport{report1}
	expectedResult = BitStatus{}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure only Field match, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test failure only Tag math - no failures should be reported
	a.ConfigFailures = []ExtendedFailure{failure0}
	a.Reports = []TestReport{report2}
	expectedResult = BitStatus{}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure only Tag match, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test failure OUT_OF range threshold - 2 reports fit rule, 1 failures should be reported with count 1 + PERCENT Exceeding type + multiple reports violation, reports meet the requirements in time range
	a.ConfigFailures = []ExtendedFailure{failure2}
	a.Reports = []TestReport{report3, report4, report5, report6}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure2.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       1,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure OUT_OF range threshold, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
		if !reflect.DeepEqual(a.Status.Failures, expectedResult.Failures) {
			t.Errorf("test failure OUT_OF range threshold, found failures diff")
		}
	}
	clearBitStatus(&a)

	// test failure VALUE Exceeding type + NO_WINDOW - multiple reports violation
	a.ConfigFailures = []ExtendedFailure{failure1}
	a.Reports = []TestReport{report3, report4, report5, report6}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure1.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       3,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure VALUE Exceeding type, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
		if !reflect.DeepEqual(a.Status.Failures, expectedResult.Failures) {
			t.Errorf("test failure VALUE Exceeding type, found failures diff")
			if len(a.Status.Failures) > 0 && a.Status.Failures[0].Count != 3 {
				t.Errorf("test failure count for NO_WINDOW, , expected 3, got %v", a.Status.Failures[0].Count)
			}
		}
	}
	clearBitStatus(&a)

	// test failure SLIDING - multiple reports violation, reports did not meet the requirements
	a.ConfigFailures = []ExtendedFailure{failure3}
	a.Reports = []TestReport{report1}
	expectedResult = BitStatus{}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure SLIDING - multiple reports violation, reports did not meet the requirements, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test failure SLIDING - multiple reports fit rule in different time frame, 1 failures should be reported with count 2
	a.ConfigFailures = []ExtendedFailure{failure3}
	a.Reports = []TestReport{report1, report6, report7, report8, report9}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure3.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       2,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure SLIDING - multiple reports violation, reports meet the requirements, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test failure SLIDING - multiple reports violation, reports meet the requirements between time ranges + trigger smaller than window
	a.ConfigFailures = []ExtendedFailure{failure2}
	a.Reports = []TestReport{report3, report5, report6}
	expectedResult = BitStatus{}
	// only 2 violations in 0-3 sec window, failure not found
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure SLIDING - multiple reports violation between time ranges, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	time.Sleep(2 * time.Second)
	// another violation in 0-3 sec window, and 3 violation between 1-4 sec window, failure count should be 2
	reports := []TestReport{report10, report9}
	a.UpdateReports(reports)
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure2.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       2,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure SLIDING - multiple reports violation between time ranges, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

}

// test for report duration, reliability of Duration method
func TestReportDuration(t *testing.T) {
	var a BitAnalyzer

	testTime := time.Now()
	expectedResult := BitStatus{}

	// test for NO_LATCH indication
	a.ConfigFailures = []ExtendedFailure{failure0}
	a.Reports = []TestReport{report0}
	// failure is reported in bit status after first crosscheck
	a.Crosscheck(testTime)
	// failure should disappear after CleanBitStatus which happens in the end of WriteBitStatus
	a.CleanBitStatus()
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for NO_LATCH indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test for LATCH_UNTIL_RESET + LATCH_FOREVER indication

	a.ConfigFailures = []ExtendedFailure{failure1}
	a.Reports = []TestReport{report3, report4, report5, report6}
	a.Crosscheck(testTime)
	// bitStatus contains 1 failure with count 3

	a.CleanBitStatus()
	a.Reports = []TestReport{}
	// bitStatus is empty, analyzer should keep UNTIL_RESET failure
	a.Crosscheck(testTime.Add(time.Second))
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure1.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       3,
			},
		},
	}
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for LATCH_UNTIL_RESET indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	// filter LATCH_UNTIL_RESET indication failures
	a.ResetSavedFailures()
	a.CleanBitStatus()

	expectedResult = BitStatus{}
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for LATCH_UNTIL_RESET indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test for LATCH_FOREVER indication
	a.ConfigFailures = []ExtendedFailure{failure2}
	a.Reports = []TestReport{report3, report4, report5, report6}
	a.Crosscheck(testTime)
	// bitStatus contains 1 failure

	a.CleanBitStatus()
	a.Reports = []TestReport{}
	// bitStatus is empty, analyzer should keep forever failure
	a.Crosscheck(testTime.Add(time.Second))
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure2.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       1,
			},
		},
	}
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for LATCH_FOREVER indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}

	clearBitStatus(&a)

	// test for NUM_OF_SECONDS indication
	a.ConfigFailures = []ExtendedFailure{failure3}
	a.Reports = []TestReport{report1, report6, report7, report8, report9}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure3.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       2,
			},
		},
	}
	a.Crosscheck(testTime)
	// bitStatus contains 1 failure with count 2
	a.CleanBitStatus()
	a.Reports = []TestReport{}
	// bitStatus is empty, analyzer should keep failure for 3 more seconds
	time.Sleep(time.Second)
	// bitStatus is empty, analyzer should keep failure for 2 more seconds
	a.Crosscheck(testTime)
	// bitStatus should contain 1 failure
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for NUM_OF_SECONDS indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}

	a.CleanBitStatus()
	a.Reports = []TestReport{}
	time.Sleep(3 * time.Second)
	a.Crosscheck(testTime)
	// bitStatus should be empty
	expectedResult = BitStatus{}
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test for NUM_OF_SECONDS indication, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)
}

// test for Dependencies, reliability of masking
func TestDependencies(t *testing.T) {
	var a BitAnalyzer

	testTime := time.Now()

	// test for masking groups, group1 should be masked, failure2 belongs to another group so it will be reported
	a.ConfigFailures = []ExtendedFailure{failure1, failure2}
	a.Reports = []TestReport{report3, report4, report5, report6}
	expectedResult := BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: failure1.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       3,
			},
			{
				FailureData: failure2.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       1,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure OUT_OF range threshold, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)

	// test for masking groups, "group1", "group temp" are masked, failure2 should not be reported
	tempFailure := failure1
	tempFailure.Failure.Dependencies.MasksOtherGroup = []string{
		"group temp",
		"group1",
	}
	a.ConfigFailures = []ExtendedFailure{tempFailure, failure2}
	a.Reports = []TestReport{report3, report4, report5, report6}
	expectedResult = BitStatus{
		Failures: []*BitStatus_RportedFailure{
			{
				FailureData: tempFailure.Failure.Description,
				Timestamp:   timestamppb.New(testTime),
				Count:       3,
			},
		},
	}
	a.Crosscheck(testTime)
	if !reflect.DeepEqual(a.Status, expectedResult) {
		t.Errorf("test failure OUT_OF range threshold, expected %v failures, got %v", len(expectedResult.Failures), len(a.Status.Failures))
	}
	clearBitStatus(&a)
}

func clearBitStatus(a *BitAnalyzer) {
	a.ConfigFailures = []ExtendedFailure{}
	a.SavedFailures = []ExtendedFailure{}
	a.Status = BitStatus{}
	a.Reports = []TestReport{}
}

// temp variables for tests

var failure0 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "volts test",
			TestId:         1,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds voltage problem",
			Purpose:        "check voltage is within 1-7 range, with a deviation of 10%",
			Severity:       1,
			OperatorFailure: []string{
				"unable to start",
				"normal functionality is damaged",
			},
			LineReplacentUnits: []string{
				"line1",
				"line2",
			},
			FieldReplacemntUnits: []string{
				"field1",
				"field2",
				"field3",
			},
		},
		ExaminationRule: &FailureExaminationRule{
			MatchingField: "volts",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname"),
				Value: []byte("server02"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       2,
					Miximum:       7,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT,
						Value: 10,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     5,
					FailuresCCount: 0,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_NO_LATCH,
			IndicationSeconds: 600,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group1",
				"groupRafael",
				"groupField",
			},
			MasksOtherGroup: []string{
				"group3",
				"group4",
			},
		},
	},
}

var failure1 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "temperature test",
			TestId:         2,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds temperature problem",
			Purpose:        "check temperature is within 60-80 range, with a deviation of 8",
			Severity:       2,
			OperatorFailure: []string{
				"can't ignite",
				"normal functionality is damaged",
			},
			LineReplacentUnits: []string{
				"line1",
				"line2",
			},
			FieldReplacemntUnits: []string{
				"field1",
				"field2",
				"field3",
			},
		},
		ExaminationRule: &FailureExaminationRule{
			MatchingField: "TemperatureCelsius",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname"),
				Value: []byte("server02"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       60,
					Miximum:       80,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_VALUE,
						Value: 8,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_NO_WINDOW,
					WindowSize:     5,
					FailuresCCount: 2,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_LATCH_UNTIL_RESET,
			IndicationSeconds: 0,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"TemperatureCelsius group",
				"groupRafael",
				"group general",
			},
			MasksOtherGroup: []string{
				"group1",
			},
		},
	},
}

var failure2 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "pressure test",
			TestId:         2,
			BitType:        []string{"CBIT, PBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds air pressure problem",
			Purpose:        "check pressure is out of 0-20 range, with a deviation of 16 percent",
			Severity:       24,
			OperatorFailure: []string{
				"OperatorFailure1",
				"OperatorFailure2",
				"OperatorFailure3",
			},
			LineReplacentUnits: []string{
				"line1",
				"line2",
			},
			FieldReplacemntUnits: []string{
				"field1",
				"field2",
				"field3",
			},
		},
		ExaminationRule: &FailureExaminationRule{
			MatchingField: "AirPressure",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname123"),
				Value: []byte("north"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       0,
					Miximum:       20,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT,
						Value: 16,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     3,
					FailuresCCount: 2,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_LATCH_FOREVER,
			IndicationSeconds: 0,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group temp",
				"group1",
			},
			MasksOtherGroup: []string{
				"group2",
			},
		},
	},
	Time: time.Now(),
}

var failure3 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "oil test",
			TestId:         1,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds oil problem",
			Purpose:        "check oil is within 0-9 range, with a deviation of 1",
			Severity:       1,
			OperatorFailure: []string{
				"unable to start",
				"normal functionality is damaged",
			},
			LineReplacentUnits: []string{
				"line1",
				"line2",
			},
			FieldReplacemntUnits: []string{
				"field1",
				"field2",
				"field3",
			},
		},
		ExaminationRule: &FailureExaminationRule{
			MatchingField: "oil",
			MatchingTag: &KeyValuePair{
				Key:   []byte("ss"),
				Value: []byte("longstanding"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       0,
					Miximum:       9,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_VALUE,
						Value: 1,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     1,
					FailuresCCount: 1,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_NUM_OF_SECONDS,
			IndicationSeconds: 3,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group1",
			},
			MasksOtherGroup: []string{
				"group4",
			},
		},
	},
}

var report0 = TestReport{
	TestId:         123,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "volts", Value: "6.5"},
	},
}

var report1 = TestReport{
	TestId:         124,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server01"},
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "volts", Value: "10"},
		{Key: "oil", Value: "4"},
	},
}

var report2 = TestReport{
	TestId:         125,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "-1"},
	},
}

var report3 = TestReport{
	TestId:         126,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "-3.3"},
		{Key: "TemperatureCelsius", Value: "68"},
	},
}

var report4 = TestReport{
	TestId:         127,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "50"},
		{Key: "TemperatureCelsius", Value: "60"},
	},
}

var report5 = TestReport{
	TestId:         129,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "10"},
		{Key: "TemperatureCelsius", Value: "72"},
	},
}

var report6 = TestReport{
	TestId:         128,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(time.Second),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "23.25"},
		{Key: "TemperatureCelsius", Value: "70"},
	},
}

var report7 = TestReport{
	TestId:         130,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "8"},
	},
}

var report8 = TestReport{
	TestId:         131,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(2 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "2"},
	},
}

var report9 = TestReport{
	TestId:         132,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(2 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "4.4"},
		{Key: "AirPressure", Value: "40.4"},
	},
}

var report10 = TestReport{
	TestId:         133,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(4 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "4.4"},
		{Key: "AirPressure", Value: "35"},
	},
}

var report11 = TestReport{
	TestId:         134,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(4 * time.Second),
}

var report12 = TestReport{
	TestId:         135,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(time.Hour),
}
