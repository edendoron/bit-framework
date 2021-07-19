package bitConfig

import (
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
)

var Configs models.ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)
}
