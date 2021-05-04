package bitHistoryCurator

import (
	"encoding/json"
	"io/ioutil"
)

func getConfig() HistoryCuratorConfig {
	configs := HistoryCuratorConfig{}
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		//TODO: handle error
	}
	err = json.Unmarshal(content, &configs)
	if err != nil {
		//TODO: handle error
	}
	return configs
}
