package bitExporter

import (
	. "../models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)

	CurrentBW.Size = Configs.BitExporterDefaultBWSize
	CurrentBW.UnitsPerSecond = Configs.BitExporterDefaultBWUnits
}
