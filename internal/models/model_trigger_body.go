package models

type TriggerBody struct {
	// Type of BIT. May be Power-On/Initiated/Continuous BIT etc.
	BitType string `json:"bitType" validate:"required"`
	// Period in [sec] resolution. 0 or negative stands for 'one-time shot'.
	PeriodSec float64 `json:"periodSec" validate:"required"`
}
