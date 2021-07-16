package bitStorageAccess

import (
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)
}
