package bitIndexer

import (
	. "../models"
	"github.com/segmentio/conf"
)

var Configs ProgConfigs

func LoadConfigs() {
	conf.Load(&Configs)

}
