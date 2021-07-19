package bitexporter

import (
	"bytes"
	"encoding/json"
	prqueue "github.com/coraxster/PriorityQueue"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	m "math"
	"net/http"
	"time"
)

// init global package variables and constants

// ReportsQueue saves the incoming reports by priority
var ReportsQueue = prqueue.Build()

// QueueChannel indicates incoming reports
var QueueChannel = make(chan int, 1000)

// CurrentBW stores current state of exporter
var CurrentBW = models.Bandwidth{}

// TestingFlag assist in tests manipulations
var TestingFlag = false

const KiB = 1024
const MiB = KiB * 1024
const GiB = MiB * 1024
const TiB = GiB * 1024
const K = 1000
const M = K * 1000
const G = M * 1000
const T = G * 1000

// ReportStatus is an enum to indicates reason for the previous epoch report sent
type ReportStatus int

const (
	SizeLimit ReportStatus = iota
	EmptyQueue
	TimeLimit
)

// ReportsScheduler manages writing reports from priority queue to storage according to CurrentBW limitations
func ReportsScheduler(duration time.Duration) {
	var indexerReqEpochSize = CalculateExtraSize()

	epoch := time.Now()
	for {
		<-QueueChannel
		for ReportsQueue.Len() > 0 {
			indexerReqEpochSize, _ = UpdateAndSendReport(indexerReqEpochSize, epoch, duration)
			if indexerReqEpochSize == -1 { // indicates that we reached size or time limit, or some error occur
				indexerReqEpochSize = CalculateExtraSize()
				epoch = time.Now()
			}
		}
	}
}

// Internal auxiliary functions

// CalculateExtraSize approximation based on wireshark results
func CalculateExtraSize() float32 {
	sizeLimitInBytes := CalculateSizeLimit(CurrentBW)
	if sizeLimitInBytes < KiB {
		return 450
	} else if sizeLimitInBytes < 4*KiB {
		return 800
	} else if sizeLimitInBytes < 100*KiB {
		return sizeLimitInBytes / 5.9
	} else if sizeLimitInBytes < 500*KiB {
		return sizeLimitInBytes / 7.8
	} else if sizeLimitInBytes < MiB {
		return sizeLimitInBytes / 12
	} else if sizeLimitInBytes < 500*MiB {
		return sizeLimitInBytes / 17
	} else {
		return sizeLimitInBytes / 25
	}
}

// CalculateSizeLimit in bytes based on @param bw
func CalculateSizeLimit(bw models.Bandwidth) float32 {
	switch bw.UnitsPerSecond {
	case "KiB":
		return bw.Size * KiB
	case "MiB":
		return bw.Size * MiB
	case "GiB":
		return bw.Size * GiB
	case "TiB":
		return bw.Size * TiB
	case "K":
		return bw.Size * K
	case "M":
		return bw.Size * M
	case "G":
		return bw.Size * G
	case "T":
		return bw.Size * T
	default:
		return 0
	}
}

// sets unlimited bandwidth size to @param bw if user requests to set negative size
func modifyBandwidthSize(bw *models.Bandwidth) {
	if bw.Size < 0 {
		bw.Size = m.MaxFloat32
	}
}

// UpdateAndSendReport function prepare and send Report to Indexer. Return the size of request body (without extra size) if queue is empty in order to know the size left to fill request, -1 otherwise. Second return value is for test purpose
func UpdateAndSendReport(epochSentSize float32, epoch time.Time, duration time.Duration) (float32, float32) {
	indexerReport := models.ReportBody{}
	postBodyPrev := &bytes.Reader{}
	for ReportsQueue.Len() > 0 && time.Since(epoch) < duration {
		item, _ := ReportsQueue.Pull()
		report := item.(models.TestReport)

		indexerReport.Reports = append(indexerReport.Reports, report)

		postBody, err := json.MarshalIndent(indexerReport, "", " ")
		if err != nil {
			log.Printf("error marshal post body to indexer, reports may be lost, error: %v", err)
			return -1, 0
		}

		sizeLimitInBytes := CalculateSizeLimit(CurrentBW)
		reportSize := epochSentSize + float32(len(postBody))

		if reportSize > sizeLimitInBytes {
			_, err = ReportsQueue.Push(report, int(report.ReportPriority))
			if err != nil {
				log.Printf("error push report to queue, reports may be too large, error: %v", err)
				return -1, 0
			}
			indexerReport.Reports = indexerReport.Reports[:len(indexerReport.Reports)-1]
			return postReportUpdateEpoch(postBodyPrev, SizeLimit, epoch, duration)
		}
		postBodyPrev = bytes.NewReader(postBody)
	}
	if time.Since(epoch) >= duration {
		return postReportUpdateEpoch(postBodyPrev, TimeLimit, epoch, duration)
	} else {
		return postReportUpdateEpoch(postBodyPrev, EmptyQueue, epoch, duration)
	}
}

/*
post reports collected so far (pulled from priority queue). return value depends on function execution reason @param status.
if due to EmptyQueue, returns to ReportsScheduler size written so far in order to keep writing new requests until size limit.
if due to SizeLimit, size limit has reached, the process will sleep for the remaining portion of the duration (1 second) and ReportsScheduler will continue to write reports in the next interval
if due to TimeLimit, time limit has reached, size sent in the current interval will be initialize and ReportsScheduler will continue to write reports.
*/
func postReportUpdateEpoch(indexerReport *bytes.Reader, status ReportStatus, epoch time.Time, duration time.Duration) (float32, float32) {
	postIndexer(indexerReport)
	switch status {
	case EmptyQueue:
		return float32(indexerReport.Size()), float32(indexerReport.Size())
	case SizeLimit:
		time.Sleep(duration - time.Since(epoch))
	}
	return -1, float32(indexerReport.Size())
}

// post reports collected under bandwidth limitations to bit-indexer service, if failed tries to post the reports back in order to insert all reports back to queue and re-send them.
func postIndexer(postBodyRef *bytes.Reader) {
	if postBodyRef.Size() < 1 || TestingFlag {
		return
	}
	go func() {
		indexerRes, err := http.Post(Configs.IndexerUrl, "application/json; charset=UTF-8", postBodyRef)
		if err != nil || indexerRes.StatusCode != http.StatusOK {
			log.Printf("error posting request to indexer, error: %v", err)
			exporterRes, e := http.Post(Configs.ExporterUrl, "application/json; charset=UTF-8", postBodyRef)
			if e != nil || exporterRes.StatusCode != http.StatusOK {
				log.Printf("error posting reports back to exporter")
				log.Fatal(e)
			}
			exporterRes.Body.Close()
			return
		}
		indexerRes.Body.Close()
	}()
}
