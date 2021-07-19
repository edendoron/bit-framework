package bitconfig

import (
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
)

// Configs is the go struct for saving all services configurations
var Configs models.ProgConfigs

// LoadConfigs extract services configurations and initialize global package variables if needed
func LoadConfigs() {
	conf.Load(&Configs)
}
