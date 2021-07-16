package bitExporter

import (
	"encoding/json"
	. "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	. "github.com/edendoron/bit-framework/internal/models"
	"net/http"
)

func GetBandwidth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(&CurrentBW)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostBandwidth(w http.ResponseWriter, r *http.Request) {
	requestBW := Bandwidth{}

	err := json.NewDecoder(r.Body).Decode(&requestBW)

	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// validate data from the user
	err = ValidateType(requestBW)
	if err != nil {
		ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// update current bandwidth
	CurrentBW = requestBW

	//validate "unitsPerSecond" units and "size" != 0
	if CalculateSizeLimit(CurrentBW) == 0 {
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
			continue
		}
		_, err = ReportsQueue.Push(report, int(report.ReportPriority))
		if err != nil {
			ApiResponseHandler(w, http.StatusBadRequest, "Bad request - report body might be too large", err)
			return
		}
	}
	QueueChannel <- ReportsQueue.Len()

	ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}
