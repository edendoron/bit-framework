package models

import (
	"time"
)

type PongBody struct {
	// Current UTC date-time (RFC3339)
	Timestamp time.Time `json:"timestamp"`
	// Service version
	Version string `json:"version"`
	// Host that running this service
	Host string `json:"host"`
	// true if operational
	Ready bool `json:"ready"`
	// API version
	ApiVersion string `json:"apiVersion"`
}
