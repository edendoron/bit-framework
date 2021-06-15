package bitHandler

import (
	. "../models"
	"fmt"
	"time"
)

var CurrentTrigger = TriggerBody{}

var TriggerChannel = make(chan TriggerBody)

var ResetIndicationChannel = make(chan bool)

const storageDataBaseUrl = "http://localhost:8082"
const storageDataReadURL = storageDataBaseUrl + "/data/read"
const storageDataWriteURL = storageDataBaseUrl + "/data/write"

func StatusScheduler() {
	d := time.Duration(CurrentTrigger.PeriodSec) * time.Second
	var analyzer BitAnalyzer
	analyzer.ReadFailuresFromStorage("config_failures")
	analyzer.ReadFailuresFromStorage("forever_failures")
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
		case <-ResetIndicationChannel:
			analyzer.FilterSavedFailures()
		case epoch := <-ticker.C:
			fmt.Println(epoch)
			go func() {
				analyzer.ReadReportsFromStorage(d)
				analyzer.Crosscheck()
				analyzer.WriteBitStatus()
			}()
		}
	}
}
