package bitExporter

import (
	exporter "github.com/edendoron/bit-framework/internal/bitTestResultsExporter"
	"github.com/edendoron/bit-framework/internal/models"
	"testing"
	"time"
)

// test PriorityQueue
func TestPriorityQueue(t *testing.T) {

	reports := []models.TestReport{report0, report1, report2}

	for _, report := range reports {
		_ = exporter.ReportsQueue.Enqueue(report, uint8(report.ReportPriority))
	}
	item, _ := exporter.ReportsQueue.Dequeue()
	report := item.(models.TestReport)
	if report.ReportPriority != 0 {
		t.Errorf("failed to prioritize reports")
	}
	item, _ = exporter.ReportsQueue.Dequeue()
	report = item.(models.TestReport)
	if report.ReportPriority != 9 {
		t.Errorf("failed to prioritize reports")
	}
	item, _ = exporter.ReportsQueue.Dequeue()
	report = item.(models.TestReport)
	if report.ReportPriority != 255 {
		t.Errorf("failed to prioritize reports")
	}

	// make sure queue can be reused
	if exporter.ReportsQueue.Len() != 0 {
		t.Errorf("Queue size expected to be %v, got %v", 0, exporter.ReportsQueue.Len())
	}

	_ = exporter.ReportsQueue.Enqueue(report2, uint8(report2.ReportPriority))
	if exporter.ReportsQueue.Len() != 1 {
		t.Errorf("Queue size expected to be %v, got %v", 1, exporter.ReportsQueue.Len())
	}
}

// test updateAndSendReport - send small report after size limit is reached
func TestUpdateAndSendReportSizeLimit(t *testing.T) {

	exporter.TestingFlag = true
	exporter.CurrentBW.Size = 0.6
	exporter.CurrentBW.UnitsPerSecond = "KiB"

	reports := []models.TestReport{report0, report1, report2}

	for _, report := range reports {
		_ = exporter.ReportsQueue.Enqueue(report, uint8(report.ReportPriority))
	}

	var indexerReqEpochSize float32 = 0
	idx := 0
	epoch := time.Now()
	for exporter.ReportsQueue.Len() > 0 {
		indexerReqEpochSize, _ = exporter.UpdateAndSendReport(indexerReqEpochSize, epoch, time.Second)
		if idx == 0 {
			if exporter.ReportsQueue.Len() != 1 {
				t.Errorf("updateAndSendReport - size limit reached, expected %v reports left in queue, got %v", 1, exporter.ReportsQueue.Len())
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

	exporter.CurrentBW.Size = 10
	exporter.CurrentBW.UnitsPerSecond = "KiB"

	reports := []models.TestReport{report0, report1, report2}

	for _, report := range reports {
		_ = exporter.ReportsQueue.Enqueue(report, uint8(report.ReportPriority))
	}

	var indexerReqEpochSize float32 = 0
	idx := 0
	epoch := time.Now()
	for exporter.ReportsQueue.Len() > 0 {
		indexerReqEpochSize, _ = exporter.UpdateAndSendReport(indexerReqEpochSize, epoch, 10*time.Nanosecond)
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
	exporter.CurrentBW.Size = 10
	exporter.CurrentBW.UnitsPerSecond = "KiB"

	reports := []models.TestReport{report0, report1, report2}
	epoch := time.Now()

	for i := 0; i < 2; i++ {
		for _, report := range reports {
			_ = exporter.ReportsQueue.Enqueue(report, uint8(report.ReportPriority))
		}

		var indexerReqEpochSize float32 = 0
		indexerReqEpochSize, _ = exporter.UpdateAndSendReport(indexerReqEpochSize, epoch, time.Second)
		if exporter.ReportsQueue.Len() != 0 {
			t.Errorf("updateAndSendReport - multiple requests, expected %v reports left in queue, got %v", 0, exporter.ReportsQueue.Len())
		}
		if indexerReqEpochSize <= 0 {
			t.Errorf("updateAndSendReport - multiple requests test failed")
		}
	}
}

var report0 = models.TestReport{
	TestId:         123,
	ReportPriority: 0,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []models.KeyValue{
		{Key: "volts", Value: "6.5"},
	},
}

var report1 = models.TestReport{
	TestId:         124,
	ReportPriority: 255,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname", Value: "server01"},
	},
	FieldSet: []models.KeyValue{
		{Key: "oil", Value: "4"},
	},
}

var report2 = models.TestReport{
	TestId:         125,
	ReportPriority: 9,
	Timestamp:      time.Now(),
	TagSet: []models.KeyValue{
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []models.KeyValue{
		{Key: "AirPressure", Value: "-1"},
	},
}
