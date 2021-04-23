package bitExporter

import (
	. "../apiResponseHandlers"
	. "../models"
	"container/heap"
	"encoding/json"
	"log"
	"net/http"
)

func GetBandwidth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// return Bandwidth response to the user
	err := json.NewEncoder(w).Encode(&currentBW)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func PostBandwidth(w http.ResponseWriter, r *http.Request) {
	request := Bandwidth{}
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal sever error", err)
		return
	}
	// validate that data from the user is of type Bandwidth
	err = ValidateType(request)
	if err != nil || calculateSizeLimit(request) == 0 {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	modifyBandwidthSize(&request)

	// update current bandwidth
	currentBW = request

	bandwidthChannel <- true

	ApiResponseHandler(w, http.StatusOK, "Bandwidth updated!", nil)
}

func ExporterPostReport(w http.ResponseWriter, r *http.Request) {
	request := ReportBody{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	// validate that data from the user is of type ReportBody
	err = ValidateType(request)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// validate the reports and insert them to the priority queue
	for _, report := range request.Reports {
		err = ValidateType(report)
		if err != nil {
			ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
			return
		}
		item := &Item{
			value:    report,
			priority: report.ReportPriority,
		}
		heap.Push(&reportsQueue, item)
	}

	// return an indication for the user that the report has received in the system
	ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}
