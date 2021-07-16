package bitHandler

import (
	. "github.com/edendoron/bit-framework/internal/models"
	"time"
)

var CurrentTrigger = TriggerBody{}

var TriggerChannel = make(chan TriggerBody)

var ResetIndicationChannel = make(chan bool)

var epoch time.Time

var status string

// StatusScheduler manages reading data from storage, analyze it and write bitStatus to storage according to CurrentTrigger
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
				status = "stopped"
			} else {
				ticker.Reset(d)
				status = "started"
			}
		case <-ResetIndicationChannel:
			analyzer.ResetSavedFailures()
			analyzer.CleanBitStatus()
		case epoch = <-ticker.C:
			go func() {
				analyzer.ReadReportsFromStorage()
				analyzer.Crosscheck(epoch)
				analyzer.WriteBitStatus()
			}()
		}
	}
}
