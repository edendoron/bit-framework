package bitHandler

import (
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)

	CurrentTrigger.PeriodSec = Configs.BitHandlerTriggerPeriod
	CurrentTrigger.BitType = Configs.BitHandlerTriggerType
}
