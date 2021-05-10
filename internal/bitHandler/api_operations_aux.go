package bitHandler

import (
	. "../apiResponseHandlers"
	. "../models"
	"net/http"
)

func handleParam(w http.ResponseWriter, r *http.Request, triggerRequest TriggerBody) bool {
	keys, ok := r.URL.Query()["action"]

	if !ok || len(keys[0]) < 1 {
		ApiResponseHandler(w, http.StatusBadRequest, "Url Param 'action' is missing", nil)
		return false
	}

	action := keys[0]
	switch action {
	case "start":
		return true
	case "stop":
		TriggerChannel <- triggerRequest
		ApiResponseHandler(w, http.StatusOK, "Trigger stopped", nil)
		return false
	default:
		ApiResponseHandler(w, http.StatusBadRequest, "Url Param 'action' is not one of 'start' or 'stop'", nil)
		return false
	}
}
