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

var queueChannel = make(chan int, 1000)

var CurrentBW = Bandwidth{}

const KiB = 1024
const MiB = KiB * 1024
const GiB = MiB * 1024
const TiB = GiB * 1024
const K = 1000
const M = K * 1000
const G = M * 1000
const T = G * 1000

const httpPostHeaderSize = 20
const reportBodyWrapSize = 218 // wireshark result
const indexerTotalExtraSize = httpPostHeaderSize + reportBodyWrapSize

type ReportStatus int

const (
	SizeLimit ReportStatus = iota
	EmptyQueue
	TimeLimit
)

// Internal auxiliary functions

func calculateSizeLimit(bw Bandwidth) float32 {
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

func ReportsScheduler(duration time.Duration) {
	var indexerReqEpochSize float32 = indexerTotalExtraSize

	epoch := time.Now()
	for {
		<-queueChannel
		for ReportsQueue.Len() > 0 {
			indexerReqEpochSize = updateReportToIndexer(indexerReqEpochSize, epoch, duration)
			if indexerReqEpochSize == -1 { // indicates that we reached size limit or time limit
				indexerReqEpochSize = indexerTotalExtraSize
				epoch = time.Now()
			}
		}
	}
}

func updateReportToIndexer(epochSentSize float32, epoch time.Time, duration time.Duration) float32 {
	indexerReport := ReportBody{}
	postBodyPrev := &bytes.Reader{}
	for ReportsQueue.Len() > 0 && time.Since(epoch) < duration {
		//TODO: handle error
		item, _ := ReportsQueue.Pull()
		report := item.(TestReport)

		indexerReport.Reports = append(indexerReport.Reports, report)

		//TODO: handle error
		postBody, _ := json.MarshalIndent(indexerReport, "", " ")

		sizeLimitInBytes := calculateSizeLimit(CurrentBW)
		reportSize := epochSentSize + float32(len(postBody))

		if reportSize > sizeLimitInBytes {
			//TODO: handle return value
			ReportsQueue.Push(report, int(report.ReportPriority))
			indexerReport.Reports = indexerReport.Reports[:len(indexerReport.Reports)-1]
			return updateRequestChannel(postBodyPrev, SizeLimit, epoch, duration)
		}
		postBodyPrev = bytes.NewReader(postBody)
	}
	if time.Since(epoch) >= duration {
		return updateRequestChannel(postBodyPrev, TimeLimit, epoch, duration)
	} else {
		return updateRequestChannel(postBodyPrev, EmptyQueue, epoch, duration)
	}
}

func updateRequestChannel(indexerReport *bytes.Reader, status ReportStatus, epoch time.Time, duration time.Duration) float32 {
	postIndexer(indexerReport)
	switch status {
	case EmptyQueue:
		return float32(indexerReport.Size())
	case SizeLimit:
		time.Sleep(duration - time.Since(epoch))
	}
	return -1
}

func postIndexer(postBodyRef *bytes.Reader) {
	if postBodyRef.Size() < 1 {
		return
	}
	go func() {
		//log.Println(time.Now(), "total size to send:", postBodyRef.Size()+indexerTotalExtraSize)
		indexerRes, err := http.Post(Configs.BitExporterPostToIndexerUrl, "application/json; charset=UTF-8", postBodyRef)
		if err != nil || indexerRes.StatusCode != http.StatusOK {
			//TODO: handle this error
			return
		}
		defer indexerRes.Body.Close()
	}()
}
