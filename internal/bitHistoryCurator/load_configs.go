package bithistorycurator

import (
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
	"log"
	"time"
)

// Configs is the go struct for saving all services configurations
var Configs models.ProgConfigs

// LoadConfigs extract services configurations and initialize global package variables if needed
func LoadConfigs() {
	conf.Load(&Configs)
}

// GetCuratorTimeConfig extract service time configuration and parse it.
func GetCuratorTimeConfig() time.Duration {
	t, err := time.ParseDuration(Configs.BitHistoryCuratorAgedDataLimit)
	if err != nil {
		log.Printf("error parsing history curator aged time duration. error: %v", err)
		return 2.628e+6
	}
	return t
}
