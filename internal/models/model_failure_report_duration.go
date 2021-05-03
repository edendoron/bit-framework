package models

type FailureReportDuration struct {
	// Specify indication type of failure
	Indication EIndicationLatchType `json:"indication" validate:"required"`
	// Relevant only if indication=NUM_OF_SECONDS
	IndicationSeconds int32 `json:"indication-seconds,omitempty"`
}

type EIndicationLatchType int

const (
	NO_LATCH          ESeverity = iota // indication is down when the test succeeds
	LATCH_UNTIL_RESET                  // indication is down until implicit request (enable delay until reported to someone)
	LATCH_FOREVER
	NUM_OF_SECONDS
)
