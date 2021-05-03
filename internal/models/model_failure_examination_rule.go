package models

type FailureExaminationRule struct {
	// The field that should be evaluated for failure report, e.g.: temperature=24.5, volts=7.1, etc.
	Field string `json:"field" validate:"required"`
	// The tag that should be presented, e.g.: hostname=server02, ip=10.1.1.1, zone=north, etc.
	Tag KeyValue `json:"tag,omitempty"`
	// The criterion that should be applied to be consider as failure
	Criteria FailureCriteria `json:"criteria" validate:"required"`
}

type FailureCriteria struct {
	// Value based definition for failure announcement
	ValueCriteria FailureValueCriteria `json:"value-criteria" validate:"required"`
	// Time based definition for failure announcement
	TimeCriteria FailureTimeCriteria `json:"time-criteria" validate:"required"`
}

type FailureValueCriteria struct {
	// minimum value of field
	Minimum float64 `json:"minimum"`
	// maximum value of field
	Maximum float64 `json:"maximum"`
	// in/out of range
	ThresholdMode EThresholdMode `json:"threshold-mode"`
	// exceeding limit - defines the allowed deviation from target threshold
	ExceedingLimit Exceeding `json:"exceeding-limit"`
}

type FailureTimeCriteria struct {
	// sliding window or no window
	WindowType   EWindowType `json:"window-type"`
	WindowSize   int32       `json:"window-size,omitempty"`
	FailureCount int32       `json:"failure-count"`
}

type Exceeding struct {
	ExceedingType  EExceedingType `json:"exceeding-type"`
	ExceedingValue float64        `json:"exceeding-value"`
}

type EThresholdMode int

const (
	IN EThresholdMode = iota
	OUT
)

type EExceedingType int

const (
	VALUE EExceedingType = iota
	PERCENT
)

type EWindowType int

const (
	NO_WINDOW EWindowType = iota
	SLIDING
)
