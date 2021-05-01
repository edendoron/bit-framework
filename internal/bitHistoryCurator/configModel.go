package bitHistoryCurator

import "time"

type HistoryCuratorConfig struct {
	AgedTimeDuration time.Duration `json:"agedTimeDuration,omitempty"`
}
