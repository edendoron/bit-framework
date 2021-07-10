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

func GetCuratorTimeConfig() time.Time {
	const layout = "2006-January-02 15:4:5"
	t, err := time.Parse(layout, Configs.BitHistoryCuratorAgedDate)
	if err != nil {
		//TODO: handle error
		return time.Now()
	}
	return t
}
