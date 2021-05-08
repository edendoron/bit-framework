package bitHandler

import (
	. "../models"
	"fmt"
	"time"
)

var CurrentTrigger = TriggerBody{}

var TriggerChannel = make(chan TriggerBody)

const storageDataBaseUrl = "http://localhost:8082"
const storageDataReadURL = storageDataBaseUrl + "/data/read"
const storageDataWriteURL = storageDataBaseUrl + "/data/write"

func StatusScheduler() {
	d := time.Duration(CurrentTrigger.PeriodSec) * time.Second
	ticker := time.NewTicker(d)
	for {
		select {
		case CurrentTrigger = <-TriggerChannel:
			d = time.Duration(CurrentTrigger.PeriodSec) * time.Second
			if d <= 0 {
				ticker.Stop()
			} else {
				ticker.Reset(d)
			}
		case epoch := <-ticker.C:
			fmt.Println(epoch)
			go func() {
				var a Analyzer
				a.ReadFromStorage("report")
				a.ReadFromStorage("failure")
				a.Crosscheck()
				a.WriteBitStatus()
			}()
		}
	}
}
