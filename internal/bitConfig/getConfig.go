package bitConfig

import (
	"encoding/json"
	"io/ioutil"
)

func getConfig() Config {
	configs := Config{}
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
