package bitHistoryCurator

import (
	. "../models"
	"github.com/segmentio/conf"
	"time"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)
}

func GetCuratorTimeConfig() time.Duration {
	t, err := time.ParseDuration(Configs.BitHistoryCuratorAgedDataLimit)
	if err != nil {
		//TODO: handle error
		return 2.628e+6
	}
	return t
}
