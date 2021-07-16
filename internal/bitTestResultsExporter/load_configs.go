package bitExporter

import (
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)

	CurrentBW.Size = Configs.BitExporterDefaultBWSize
	CurrentBW.UnitsPerSecond = Configs.BitExporterDefaultBWUnits
}
