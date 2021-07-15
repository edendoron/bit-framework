package bitHandler

import (
	. "../apiResponseHandlers"
	. "../models"
	"errors"
	"net/http"
)

// function handles Trigger according to "action" parameter.
// function returns true and stops bitStatus routine if "action" param is "stop",
// returns false if action param is "start", and returns true with proper error message to user otherwise.
func handleNonStartAction(w http.ResponseWriter, r *http.Request, triggerRequest TriggerBody) bool {
	keys, ok := r.URL.Query()["action"]

	if !ok || len(keys[0]) < 1 {
		ApiResponseHandler(w, http.StatusBadRequest, "Url Param 'action' is missing", nil)
		return true
	}

	action := keys[0]
	switch action {
	case "start":
		return false
	case "stop":
		TriggerChannel <- triggerRequest
		ApiResponseHandler(w, http.StatusOK, "Trigger stopped", nil)
		return true
	default:
		ApiResponseHandler(w, http.StatusBadRequest, "Url Param 'action' is not one of 'start' or 'stop'", errors.New("bad request"))
		return true
	}
}
