package bitExporter

import (
	. "../models"
	"bytes"
	"encoding/json"
	"fmt"
	m "math"
	"net/http"
	"time"
)

// init global package variables
var reportsQueue = make(PriorityQueue, 0)

var currentBW = Bandwidth{
	Size:           25,
	UnitsPerSecond: "KiB",
}

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

func reportsScheduler(d time.Duration, reporter func() float32) {
	for x := range time.Tick(d) {
		// calculate current bandwidth limitation
		sizeLimitInBytes := calculateSizeLimit(currentBW)
		var sizeSentInBytes float32 = 0
		fmt.Println(x)
		// sent a request to indexer service according to system current bandwidth limitation
		// NOTE: reporter may exceed the limitation if a post request already initiated.
		// TODO: need to check average postIndexer running time to determine possible average exceeding size
		for validateBandwidthLimit(sizeSentInBytes, sizeLimitInBytes) && time.Since(x) < d {
			sizeSentInBytes += reporter()
			fmt.Println("sent and then limit ", sizeSentInBytes, sizeLimitInBytes)
		}
	}
}

func validateBandwidthLimit(sizeSentInBytes float32, sizeLimitInBytes float32) bool {
	if reportsQueue.Len() < 1 {
		return false
	}
	item := reportsQueue.Top().(*Item)
	report := item.value
	postBody, _ := json.MarshalIndent(report, "", " ")
	return sizeSentInBytes+float32(len(postBody)) <= sizeLimitInBytes
}

func postIndexer() float32 {
	item := reportsQueue.Pop().(*Item)
	report := item.value
	postBody, err := json.MarshalIndent(report, "", " ")
	postBodyRef := bytes.NewBuffer(postBody)
	indexerRes, err := http.Post(postIndexedUrl, "application/json; charset=UTF-8", postBodyRef)
	if err != nil || indexerRes.StatusCode != http.StatusOK {
		//TODO: handle error on inter-services errors
		return 0
	}
	//NOTE: we return the size of postBody (actual report) and not the whole request (which includes more properties)
	return float32(len(postBody))
}
