package bitQuery

import (
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	query "github.com/edendoron/bit-framework/internal/bitQuery"
	"github.com/edendoron/bit-framework/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
	"time"
)

// test filterBitStatus
func TestFilterBitStatus1BitStatus(t *testing.T) {

	// filter 1 failure from 1 bitStatus
	var statusSlice = []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure0, failure1}},
	}

	var expectedResult = []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure1}},
	}
	maskedTests := []uint64{1}

	query.FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}
}

func TestFilterBitStatus2BitStatus(t *testing.T) {
	// filter 2 failures from different bitStatus
	statusSlice := []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure0, failure1}},
		{Failures: []*bit.BitStatus_RportedFailure{failure2, failure3}},
	}

	expectedResult := []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure1}},
		{Failures: []*bit.BitStatus_RportedFailure{failure2}},
	}
	maskedTests := []uint64{1, 4}

	query.FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}
}

func TestFilterBitStatusAllFailures(t *testing.T) {
	// filter all failures in one bit status (bit status array is not empty)
	statusSlice := []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure0, failure1}},
		{Failures: []*bit.BitStatus_RportedFailure{failure2, failure3}},
	}

	expectedResult := []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure2, failure3}},
	}
	maskedTests := []uint64{1, 2}

	query.FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}

	// filter all failures (bit status array is empty)
	statusSlice = []bit.BitStatus{
		{Failures: []*bit.BitStatus_RportedFailure{failure0, failure1}},
	}

	expectedResult = []bit.BitStatus{}
	maskedTests = []uint64{1, 2}

	query.FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}
}

// test filterReports
func TestFilterReports(t *testing.T) {

	// filter 2 reports by tag
	var reportSlice = []models.TestReport{report0, report1, report2}

	var expectedResult = []models.TestReport{report2}

	filter := "tag"
	var values = []string{"hostname123", "north"}

	query.FilterReports(&reportSlice, filter, values)

	if !reflect.DeepEqual(reportSlice, expectedResult) {
		if len(expectedResult) != len(reportSlice) {
			t.Errorf("test failure filterReports, expected %v reports, got %v", len(expectedResult), len(reportSlice))
			return
		}
	}

	// filter all reports by field
	reportSlice = []models.TestReport{report0, report1, report2}

	expectedResult = []models.TestReport{}

	filter = "field"
	values = []string{"temp"}

	query.FilterReports(&reportSlice, filter, values)

	if !reflect.DeepEqual(reportSlice, expectedResult) {
		if len(expectedResult) != len(reportSlice) {
			t.Errorf("test failure filterReports, expected %v reports, got %v", len(expectedResult), len(reportSlice))
			return
		}
	}

	// filter no reports by tag
	reportSlice = []models.TestReport{report3, report4, report5}

	expectedResult = []models.TestReport{report3, report4, report5}

	filter = "tag"
	values = []string{"hostname", "server02"}

	query.FilterReports(&reportSlice, filter, values)

	if !reflect.DeepEqual(reportSlice, expectedResult) {
		if len(expectedResult) != len(reportSlice) {
			t.Errorf("test failure filterReports, expected %v reports, got %v", len(expectedResult), len(reportSlice))
			return
		}
	}
}

// temp variables for tests

var failure0 = &bit.BitStatus_RportedFailure{
	FailureData: &bit.FailureDescription{
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
}

var failure1 = &bit.BitStatus_RportedFailure{
	FailureData: &bit.FailureDescription{
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
}

var failure2 = &bit.BitStatus_RportedFailure{
	FailureData: &bit.FailureDescription{
		UnitName:       "system test check",
		TestName:       "pressure test",
		TestId:         3,
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
	Timestamp: timestamppb.Now(),
}

var failure3 = &bit.BitStatus_RportedFailure{
	FailureData: &bit.FailureDescription{
		UnitName:       "system test check",
		TestName:       "oil test",
		TestId:         4,
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
}

var report0 = models.TestReport{
	TestId:         123,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []models.KeyValue{
		{Key: "volts", Value: "6.5"},
	},
}

var report1 = models.TestReport{
	TestId:         124,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname", Value: "server01"},
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []models.KeyValue{
		{Key: "volts", Value: "10"},
		{Key: "oil", Value: "4"},
	},
}

var report2 = models.TestReport{
	TestId:         125,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []models.KeyValue{
		{Key: "AirPressure", Value: "-1"},
	},
}

var report3 = models.TestReport{
	TestId:         126,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []models.KeyValue{
		{Key: "AirPressure", Value: "-3.3"},
		{Key: "TemperatureCelsius", Value: "68"},
	},
}

var report4 = models.TestReport{
	TestId:         127,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []models.KeyValue{
		{Key: "AirPressure", Value: "50"},
		{Key: "TemperatureCelsius", Value: "60"},
	},
}

var report5 = models.TestReport{
	TestId:         129,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []models.KeyValue{
		{Key: "AirPressure", Value: "10"},
		{Key: "TemperatureCelsius", Value: "72"},
	},
}
