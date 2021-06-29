package bitHandler

import (
	. "../models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)

	CurrentTrigger.PeriodSec = Configs.BitHandlerTriggerPeriod
	CurrentTrigger.BitType = Configs.BitHandlerTriggerType
}
