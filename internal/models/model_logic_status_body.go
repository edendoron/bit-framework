package models

import (
	"time"
)

type LogicStatusBody struct {
	Trigger *TriggerBody `json:"trigger" validate:"required"`
	Status  string       `json:"status" validate:"required"`
	// Last time BIT logic was initiated as UTC date-time (RFC3339)
	LastBitStartTimestamp time.Time `json:"lastBitStartTimestamp,omitempty" validate:"required"`
}
