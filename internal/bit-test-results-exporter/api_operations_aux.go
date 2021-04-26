package bitExporter

import (
	. "../models"
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	m "math"
	"net/http"
	"time"
)

// init global package variables
var reportsQueue = make(PriorityQueue, 0)

//var readyToPostChannel = make(chan bool)

var reportReadyChannel = make(chan bool)

//var timeLimitReadyChannel = make(chan bool)

var reportSentSizeChannel = make(chan float32)

var currentBW = Bandwidth{
	Size:           0.5,
	UnitsPerSecond: "K",
}

var indexerReport = ReportBody{}

const KiB = 1024
const MiB = KiB * 1024
const GiB = MiB * 1024
const TiB = GiB * 1024
const K = 1000
const M = K * 1000
const G = M * 1000
const T = G * 1000

const postIndexedUrl = "http://localhost:8081/report/raw"

// Internal auxiliary functions
//func writeBandwidth(request Bandwidth) ApiResponse {
//	input, err := json.MarshalIndent(request, "", " ")
//	if err != nil {
//		return ApiResponse{Code: 404, Message: "Bad request"}
//	}
//	err = ioutil.WriteFile("storage/test.json", input, 0644)
//	if err != nil {
//		log.Println(err)
//		return ApiResponse{Code: 404, Message: "Corrupt file"}
//	}
//	return ApiResponse{Code: 200, Message: "Bandwidth updated!"}
//}
//
//func readBandwidth() ApiResponse {
//	content, err := ioutil.ReadFile("storage/test.json")
//	if err != nil {
//		return ApiResponse{Code: 404, Message: "Corrupt file"}
//	}
//	return ApiResponse{Code: 200, Message: string(content)}
//}
//
//func readValidatedBandwidth(w http.ResponseWriter) (Bandwidth, bool) {
//	bw := readBandwidth()
//	if bw.Code != 200 {
//		w.WriteHeader(http.StatusNotFound)
//		err := json.NewEncoder(w).Encode(&bw)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//		}
//		return Bandwidth{}, true
//	}
//	response := Bandwidth{}
//	err := json.Unmarshal([]byte(bw.Message), &response)
//
//	// validate that data from file is ok
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		err = json.NewEncoder(w).Encode(&bw)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//		}
//		return Bandwidth{}, true
//	}
//
//	// validate that data from file is of type Bandwidth
//	if !ValidateType(response) {
//		w.WriteHeader(http.StatusNotFound)
//		err = json.NewEncoder(w).Encode(ApiResponse{Code: 404, Message: "Error reading file"})
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			log.Fatalln(err)
//		}
//		return Bandwidth{}, true
//	}
//	return response, false
//}

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

func updateReportToIndexer() {
	var reportSize float32 = 0
	for {
		if reportsQueue.Len() == 0 && len(indexerReport.Reports) > 0 {
			reportReadyChannel <- true
			reportSize = <-reportSentSizeChannel
		} else if reportsQueue.Len() > 0 {
			item := reportsQueue.Pop().(*Item)
			report := item.value
			postBody, _ := json.MarshalIndent(report, "", " ")
			sizeLimitInBytes := calculateSizeLimit(currentBW)
			if reportSize+float32(len(postBody)) <= sizeLimitInBytes {
				indexerReport.Reports = append(indexerReport.Reports, report)
				reportSize += float32(len(postBody))
				fmt.Println("added to indexerReport so far and then limit ", reportSize, sizeLimitInBytes)
			} else {
				heap.Push(&reportsQueue, item)
				reportReadyChannel <- true
				reportSize = <-reportSentSizeChannel
			}
		}
	}
}

func reportsScheduler(d time.Duration) {
	for {
		select {
		case res := <-time.After(d):
			fmt.Println("Current time", res)
			if len(indexerReport.Reports) > 0 {
				//readyToPostChannel <- true
				postIndexer()
				reportSentSizeChannel <- 0
			}
		case res := <-reportReadyChannel:
			fmt.Println("Report ready", res)
			//readyToPostChannel <- true
			reportSentSizeChannel <- postIndexer()
		}
	}
}

func postIndexer() float32 {
	//for {
	//<-readyToPostChannel
	if len(indexerReport.Reports) == 0 {
		return 0
	}
	postBody, _ := json.MarshalIndent(indexerReport, "", " ")

	//clear indexerReport
	indexerReport = ReportBody{}
	//reportSentSizeChannel <- float32(len(postBody))

	postBodyRef := bytes.NewBuffer(postBody)
	indexerRes, err := http.Post(postIndexedUrl, "application/json; charset=UTF-8", postBodyRef)
	if err != nil || indexerRes.StatusCode != http.StatusOK {
		//TODO: handle this error
		return 0
	}
	fmt.Println("sent postBody len", float32(len(postBody)))
	return float32(len(postBody))
	//}
}

//
//func updateReportToIndexer2() {
//	var reportSize float32 = 0
//	sizeLimitInBytes := calculateSizeLimit(currentBW)
//	for reportsQueue.Len() > 0 {
//		item := reportsQueue.Pop().(*Item)
//		report := item.value
//		postBody, _ := json.MarshalIndent(report, "", " ")
//		if reportSize+float32(len(postBody)) <= sizeLimitInBytes {
//			indexerReport.Reports = append(indexerReport.Reports, report)
//			reportSize += float32(len(postBody))
//			fmt.Println("added to indexerReport so far and then limit ", reportSize, sizeLimitInBytes)
//		} else {
//			heap.Push(&reportsQueue, item)
//			return
//		}
//		indexerReport.Reports = append(indexerReport.Reports, report)
//	}
//}
//
//func postIndexer2() {
//	if len(indexerReport.Reports) == 0 {
//		return
//	}
//	postBody, _ := json.MarshalIndent(indexerReport, "", " ")
//	src := bytes.NewReader(postBody)
//	//src := bytes.NewReader(make([]byte, 2687))
//	rate := calculateSizeLimit(currentBW)
//	fmt.Println(rate)
//	bucket := ratelimit.NewBucketWithRate(float64(rate), int64(rate))
//
//	//dst := &bytes.Buffer{}
//
//	start := time.Now()
//
//	// Copy source to destination, but wrap our reader with rate limited one
//	//io.Copy(dst, ratelimit.Reader(src, bucket))
//	client := &http.Client{}
//	req, _ := http.NewRequest("POST", postIndexedUrl, ratelimit.Reader(src, bucket))
//	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
//	resp, _ := client.Do(req)
//	fmt.Println("resp code", resp.StatusCode)
//
//	//indexerRes, _ := http.Post(postIndexedUrl, "application/json; charset=UTF-8", ratelimit.Reader(src, bucket))
//	b, _ := ioutil.ReadAll(req.Body)
//	fmt.Printf("Sent %d bytes in %s\n", len(b), time.Since(start))
//	//if err != nil || indexerRes.StatusCode != http.StatusOK {
//	//	//TODO: handle this error
//	//	return
//	//}
//}
