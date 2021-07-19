package bitexporter

import (
	"encoding/json"
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	"github.com/edendoron/bit-framework/internal/models"
	"net/http"
)

func GetBandwidth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(&CurrentBW)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostBandwidth(w http.ResponseWriter, r *http.Request) {
	requestBW := models.Bandwidth{}

	err := json.NewDecoder(r.Body).Decode(&requestBW)

	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// validate data from the user
	err = models.ValidateType(requestBW)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// update current bandwidth
	CurrentBW = requestBW

	//validate "unitsPerSecond" units and "size" != 0
	if CalculateSizeLimit(CurrentBW) == 0 {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	modifyBandwidthSize(&CurrentBW)

	// return ApiResponse response to the user
	rh.ApiResponseHandler(w, http.StatusOK, "Bandwidth updated!", nil)
}

func ExporterPostReport(w http.ResponseWriter, r *http.Request) {
	request := models.ReportBody{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	// validate that data from the user is of type ReportBody
	err = models.ValidateType(request)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	// validate the reports and insert them to the priority queue
	for _, report := range request.Reports {
		err = models.ValidateType(report)
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request", err)
			continue
		}
		_, err = ReportsQueue.Push(report, int(report.ReportPriority))
		if err != nil {
			rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad request - report body might be too large", err)
			return
		}
	}
	QueueChannel <- ReportsQueue.Len()

	rh.ApiResponseHandler(w, http.StatusOK, "Report received!", nil)
}
