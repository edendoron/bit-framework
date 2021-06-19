package models

import (
	. "../../configs/rafael.com/bina/bit"
	"encoding/json"
)

type BitStatusList struct {
	// Multiple bit status reports set
	StatusList []BitStatus `json:"status_list" validate:"required"`
}

func (bsl *BitStatusList) String() string {
	jsonRes, err := json.Marshal(bsl)
	if err != nil {
		return ""
	}
	return string(jsonRes)
}
