package bitHistoryCurator

import (
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
	"log"
	"time"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)
}

func GetCuratorTimeConfig() time.Duration {
	t, err := time.ParseDuration(Configs.BitHistoryCuratorAgedDataLimit)
	if err != nil {
		log.Printf("error parsing history curator aged time duration. error: %v", err)
		return 2.628e+6
	}
	return t
}
