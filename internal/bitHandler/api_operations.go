package bithandler

import (
	"encoding/json"
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	"github.com/edendoron/bit-framework/internal/models"
	"net/http"
)

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	result := models.LogicStatusBody{
		Trigger:               &CurrentTrigger,
		LastBitStartTimestamp: epoch,
		Status:                status,
	}

	err := json.NewEncoder(w).Encode(&result)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	triggerRequest := models.TriggerBody{}

	if handleNonStartAction(w, r, triggerRequest) {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&triggerRequest)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad Request", err)
		return
	}

	// validate that data from the user is of type TriggerBody
	err = models.ValidateType(triggerRequest)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad Request", err)
		return
	}

	// update current TriggerBody
	TriggerChannel <- triggerRequest

	// return ApiResponse response to the user
	rh.ApiResponseHandler(w, http.StatusOK, "Trigger updated!", nil)
	w.WriteHeader(http.StatusOK)
}

func PutResetIndications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// remove irrelevant indications
	ResetIndicationChannel <- true

	// return ApiResponse response to the user
	rh.ApiResponseHandler(w, http.StatusOK, "Reset indications!", nil)
	w.WriteHeader(http.StatusOK)
}
