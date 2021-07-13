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
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	if code != 200 {
		log.Println("ApiResponseHandler error is: ", e)
	}
}
