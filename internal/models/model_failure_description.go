package models

type FailureDescription struct {
	// Unit name that performs the test
	UnitName string `json:"unit-name" validate:"required"`
	// Test name
	TestName string `json:"test-name" validate:"required"`
	// Test ID, must be unique in the system
	TestId int64 `json:"test-id" validate:"required"`
	// Type of BIT, e.g. PBIT, IBIT, CBIT etc
	BitType []string `json:"bit-type" validate:"required"`
	// Test procedure description for presentation
	Description string `json:"description" validate:"required"`
	// Test additional information
	AdditionalInfo string `json:"additional-info" validate:"required"`
	// Test purpose
	Purpose string `json:"purpose" validate:"required"`
	// Severity of test failure
	Severity ESeverity `json:"severity" validate:"required"`
	// This Failure should not be reported to user groups
	DoNotReportUserGroup []string `json:"do-not-report-user-group" validate:"required"`
	// Describes system functionalities from operators point of view that are influenced (what we cannot do in the system) in case of failure e.g. TargetDetection
	OperatorFailure []string `json:"operator-failure" validate:"required"`
	// Line replacement units, in the order of replacement
	LineReplacementUnits []string `json:"line-replacement-units" validate:"required"`
	// Field replacement units, in the order of replacement
	FieldReplacementUnits []string `json:"field-replacement-units" validate:"required"`
}

type ESeverity int

const (
	MINOR ESeverity = iota
	DEGRADED
	CRITICAL
	SAFETY
)
