package bitExporter

import (
	. "../models"
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// init global package variables
var reportsQueue = make(PriorityQueue, 0)

var currentBW = Bandwidth{
	Size:           25,
	UnitsPerSecond: "MB",
}

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
		return bw.Size * 1000
	case "MB":
		return bw.Size * 1e6
	default:
		return bw.Size * 1000
	}
}

func reportsScheduler(d time.Duration, reporter func() float32) {
	for x := range time.Tick(d) {
		// calculate current bandwidth limitation
		sizeLimitInBytes := calculateSizeLimit(currentBW)
		var sizeSentInBytes float32 = 0
		// fmt.Println(x)
		// sent a request to indexer service according to system current bandwidth limitation
		// NOTE: reporter may exceed the limitation if a post request already initiated.
		// TODO: need to check average postIndexer running time to determine possible average exceeding size
		for reportsQueue.Len() > 0 && sizeSentInBytes < sizeLimitInBytes && time.Since(x) < d {
			sizeSentInBytes += reporter()
			fmt.Println("sent and then limit ", sizeSentInBytes, sizeLimitInBytes)
		}
	}
}

func postIndexer() float32 {
	item := heap.Pop(&reportsQueue).(*Item)
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
