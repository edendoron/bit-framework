package bitQuery

import (
	. ".."
	. "../../../configs/rafael.com/bina/bit"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
)

// test filterBitStatus
func TestFilterBitStatus(t *testing.T) {

	// filter 1 failure from 1 bitStatus
	var statusSlice = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure0, failure1}},
	}

	var expectedResult = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure1}},
	}
	maskedTests := []uint64{1}

	FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i, _ := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}

	// filter 2 failures from different bitStatus
	statusSlice = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure0, failure1}},
		{Failures: []*BitStatus_RportedFailure{failure2, failure3}},
	}

	expectedResult = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure1}},
		{Failures: []*BitStatus_RportedFailure{failure2}},
	}
	maskedTests = []uint64{1, 4}

	FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i, _ := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}

	// filter all failures in one bit status (bit status array is not empty)
	statusSlice = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure0, failure1}},
		{Failures: []*BitStatus_RportedFailure{failure2, failure3}},
	}

	expectedResult = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure2, failure3}},
	}
	maskedTests = []uint64{1, 2}

	FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i, _ := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}

	// filter all failures (bit status array is empty)
	statusSlice = []BitStatus{
		{Failures: []*BitStatus_RportedFailure{failure0, failure1}},
	}

	expectedResult = []BitStatus{}
	maskedTests = []uint64{1, 2}

	FilterBitStatus(&statusSlice, maskedTests)

	if !reflect.DeepEqual(statusSlice, expectedResult) {
		if len(expectedResult) != len(statusSlice) {
			t.Errorf("test failure filterBitStatus, expected %v statuses, got %v", len(expectedResult), len(statusSlice))
			return
		}
		for i, _ := range statusSlice {
			if !reflect.DeepEqual(expectedResult[i].Failures, statusSlice[i].Failures) {
				t.Errorf("test failure filterBitStatus, status array at index %v, expected %v failures, got %v", i, len(expectedResult[i].Failures), len(statusSlice[i].Failures))
			}
		}
	}
}

// temp variables for tests

var failure0 = &BitStatus_RportedFailure{
	FailureData: &FailureDescription{
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

var failure1 = &BitStatus_RportedFailure{
	FailureData: &FailureDescription{
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

var failure2 = &BitStatus_RportedFailure{
	FailureData: &FailureDescription{
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

var failure3 = &BitStatus_RportedFailure{
	FailureData: &FailureDescription{
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
