package bitExporter

import (
	. ".."
	. "../../models"
	"testing"
	"time"
)

// test PriorityQueue
func TestPriorityQueue(t *testing.T) {

	reports := []TestReport{report0, report1, report2}

	for _, report := range reports {
		push, err := ReportsQueue.Push(report, int(report.ReportPriority))
		if push == false || err != nil {
			t.Errorf("Push reports to queue failed %v", err)
			return
		}
	}
	item, _ := ReportsQueue.Pull()
	report := item.(TestReport)
	if report.ReportPriority != 1 {
		t.Errorf("failed to prioritize reports")
	}
	item, _ = ReportsQueue.Pull()
	report = item.(TestReport)
	if report.ReportPriority != 9 {
		t.Errorf("failed to prioritize reports")
	}
	item, _ = ReportsQueue.Pull()
	report = item.(TestReport)
	if report.ReportPriority != 120 {
		t.Errorf("failed to prioritize reports")
	}

}

// test updateAndSendReport - send small report after size limit is reached
func TestUpdateAndSendReportSizeLimit(t *testing.T) {

	CurrentBW.Size = 0.6
	CurrentBW.UnitsPerSecond = "KiB"

	reports := []TestReport{report0, report1, report2}

	for _, report := range reports {
		push, err := ReportsQueue.Push(report, int(report.ReportPriority))
		if push == false || err != nil {
			t.Errorf("Push reports to queue failed %v", err)
			return
		}
	}

	var indexerReqEpochSize float32 = 0
	idx := 0
	epoch := time.Now()
	for ReportsQueue.Len() > 0 {
		indexerReqEpochSize, _ = UpdateAndSendReport(indexerReqEpochSize, epoch, time.Second)
		if idx == 0 {
			if ReportsQueue.Len() != 1 {
				t.Errorf("updateAndSendReport - size limit reached, expected %v reports left in queue, got %v", 1, ReportsQueue.Len())
			}
			if indexerReqEpochSize != -1 {
				t.Errorf("updateAndSendReport - size limit reached test failed")
			} else {
				indexerReqEpochSize = 0
				epoch = time.Now()
			}
		}
		if idx == 1 {
			if indexerReqEpochSize <= 0 {
				t.Errorf("updateAndSendReport - send small report after size limit is reached test failed")
			}
		}

		idx++
	}
}

// test updateAndSendReport - send report when time limit reached
func TestUpdateAndSendReportTimeLimit(t *testing.T) {

	CurrentBW.Size = 10
	CurrentBW.UnitsPerSecond = "KiB"

	reports := []TestReport{report0, report1, report2}

	for _, report := range reports {
		push, err := ReportsQueue.Push(report, int(report.ReportPriority))
		if push == false || err != nil {
			t.Errorf("Push reports to queue failed %v", err)
			return
		}
	}

	var indexerReqEpochSize float32 = 0
	idx := 0
	epoch := time.Now()
	for ReportsQueue.Len() > 0 {
		indexerReqEpochSize, _ = UpdateAndSendReport(indexerReqEpochSize, epoch, 10*time.Nanosecond)
		if idx == 0 {
			if indexerReqEpochSize != -1 {
				t.Errorf("updateAndSendReport - time limit reached test failed")
			} else {
				indexerReqEpochSize = 0
				epoch = time.Now()
			}
		}

		idx++
	}
}

// test updateAndSendReport - send request multiple times in timeframe
func TestUpdateAndSendReportMultipleReports(t *testing.T) {
	CurrentBW.Size = 10
	CurrentBW.UnitsPerSecond = "KiB"

	reports := []TestReport{report0, report1, report2}
	epoch := time.Now()

	for i := 0; i < 2; i++ {
		for _, report := range reports {
			push, err := ReportsQueue.Push(report, int(report.ReportPriority))
			if push == false || err != nil {
				t.Errorf("Push reports to queue failed %v", err)
				return
			}
		}

		var indexerReqEpochSize float32 = 0
		indexerReqEpochSize, _ = UpdateAndSendReport(indexerReqEpochSize, epoch, time.Second)
		if ReportsQueue.Len() != 0 {
			t.Errorf("updateAndSendReport - multiple requests, expected %v reports left in queue, got %v", 0, ReportsQueue.Len())
		}
		if indexerReqEpochSize <= 0 {
			t.Errorf("updateAndSendReport - multiple requests test failed")
		}
	}
}

var report0 = TestReport{
	TestId:         123,
	ReportPriority: 1,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "volts", Value: "6.5"},
	},
}

var report1 = TestReport{
	TestId:         124,
	ReportPriority: 120,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server01"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "4"},
	},
}

var report2 = TestReport{
	TestId:         125,
	ReportPriority: 9,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "-1"},
	},
}
