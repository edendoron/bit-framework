package bitHandler

import (
	. "../models"
	"time"
)

var CurrentTrigger = TriggerBody{}

var TriggerChannel = make(chan TriggerBody)

var ResetIndicationChannel = make(chan bool)

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
			analyzer.ResetSavedFailures()
			analyzer.CleanBitStatus()
		case epoch := <-ticker.C:
			//fmt.Println(epoch)
			go func() {
				analyzer.ReadReportsFromStorage(d)
				analyzer.Crosscheck(epoch)
				analyzer.WriteBitStatus()
			}()
		}
	}
}
