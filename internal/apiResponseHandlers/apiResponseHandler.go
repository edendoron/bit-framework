package apiResponseHandlers

import (
	. "../models"
	"encoding/json"
	"log"
	"net/http"
)

func ApiResponseHandler(w http.ResponseWriter, code int, message string, e error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := ApiResponse{Code: int32(code), Message: message}
	w.WriteHeader(code)
	if code != 200 {
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	}
}
