package bitExporter

import (
	. "../apiResponseHandlers"
	. "../models"
	"encoding/json"
	"log"
	"net/http"
)

func GetBandwidth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// return Bandwidth response to the user
	err := json.NewEncoder(w).Encode(&CurrentBW)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}
	w.WriteHeader(http.StatusOK)
}

func PostBandwidth(w http.ResponseWriter, r *http.Request) {
	request := Bandwidth{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	// validate that data from the user is of type Bandwidth
	err = ValidateType(request)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// update current bandwidth
	CurrentBW = request

	//validate "unitsPerSecond" units and "size" > 0
	if calculateSizeLimit(CurrentBW) == 0 {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	modifyBandwidthSize(&CurrentBW)

	// return ApiResponse response to the user
	ApiResponseHandler(w, http.StatusOK, "Bandwidth updated!", nil)
}

func ExporterPostReport(w http.ResponseWriter, r *http.Request) {
	request := ReportBody{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
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
		// TODO: handle error
		ReportsQueue.Push(report, int(report.ReportPriority))
	}
	queueChannel <- ReportsQueue.Len()

	ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}
