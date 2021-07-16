package models

import (
	"time"
)

type TestReport struct {
	// The unique ID of the performed test
	TestId float64 `json:"testId" validate:"required"`
	// The report priority [0 - lowest]
	ReportPriority float64 `json:"reportPriority,omitempty" validate:"required"`
	// UTC date-time (RFC3339) when the test was performed
	Timestamp time.Time `json:"timestamp" validate:"required"`
	// set of Tags, e.g.: hostname=server02, ip=10.1.1.1, zone=north, etc.
	TagSet []KeyValue `json:"tagSet,omitempty" validate:"required"`
	// set of Fields, e.g.: temperature=24.5, volts=7.1, etc.
	FieldSet []KeyValue `json:"fieldSet,omitempty" validate:"required"`
}
