package bitExporter

import (
	. "../models"
	"bytes"
	"encoding/json"
	prqueue "github.com/coraxster/PriorityQueue"
	m "math"
	"net/http"
	"time"
)

// init global package variables and constants

var ReportsQueue = prqueue.Build()

var QueueChannel = make(chan int, 1000)

var CurrentBW = Bandwidth{}

const KiB = 1024
const MiB = KiB * 1024
const GiB = MiB * 1024
const TiB = GiB * 1024
const K = 1000
const M = K * 1000
const G = M * 1000
const T = G * 1000

type ReportStatus int

const (
	SizeLimit ReportStatus = iota
	EmptyQueue
	TimeLimit
)

func ReportsScheduler(duration time.Duration) {
	var indexerReqEpochSize = CalculateExtraSize()

	epoch := time.Now()
	for {
		<-QueueChannel
		for ReportsQueue.Len() > 0 {
			indexerReqEpochSize, _ = UpdateAndSendReport(indexerReqEpochSize, epoch, duration)
			if indexerReqEpochSize == -1 { // indicates that we reached size limit or time limit
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

func CalculateSizeLimit(bw Bandwidth) float32 {
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

func modifyBandwidthSize(bw *Bandwidth) {
	if bw.Size <= 0 {
		bw.Size = m.MaxFloat32
	}
}

// UpdateAndSendReport function prepare and send Report to Indexer. Return the size of request body (without extra size) if queue is empty in order to know the size left to fill request, -1 otherwise. Second return value is for test purpose
func UpdateAndSendReport(epochSentSize float32, epoch time.Time, duration time.Duration) (float32, float32) {
	indexerReport := ReportBody{}
	postBodyPrev := &bytes.Reader{}
	for ReportsQueue.Len() > 0 && time.Since(epoch) < duration {
		//TODO: handle error
		item, _ := ReportsQueue.Pull()
		report := item.(TestReport)

		indexerReport.Reports = append(indexerReport.Reports, report)

		//TODO: handle error
		postBody, _ := json.MarshalIndent(indexerReport, "", " ")

		sizeLimitInBytes := CalculateSizeLimit(CurrentBW)
		reportSize := epochSentSize + float32(len(postBody))

		if reportSize > sizeLimitInBytes {
			//TODO: handle return value
			ReportsQueue.Push(report, int(report.ReportPriority))
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

func postIndexer(postBodyRef *bytes.Reader) {
	if postBodyRef.Size() < 1 {
		return
	}
	go func() {
		//fmt.Println(time.Now(), "total size to send:", postBodyRef.Size())
		//fmt.Println(time.Now(), "total size to send + extra size:", float32(postBodyRef.Size())+CalculateExtraSize())
		indexerRes, err := http.Post(Configs.BitExporterPostToIndexerUrl, "application/json; charset=UTF-8", postBodyRef)
		if err != nil || indexerRes.StatusCode != http.StatusOK {
			//TODO: handle this error
			return
		}
		defer indexerRes.Body.Close()
	}()
}
