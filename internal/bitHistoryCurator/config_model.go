package bitHistoryCurator

import "time"

type HistoryCuratorConfig struct {
	AgedTimeDuration time.Time `json:"agedTimeDuration,omitempty"`
}
