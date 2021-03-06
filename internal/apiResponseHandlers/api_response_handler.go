package apiResponseHandlers

import (
	"encoding/json"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
)

// ApiResponseHandler general purpose response handler that writes 'ApiResponse' to @w
func ApiResponseHandler(w http.ResponseWriter, code int, message string, e error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := models.ApiResponse{Code: int32(code), Message: message}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	if code != 200 {
		log.Printf("%v \n ApiResponseHandler error is: %v", message, e)
	}
}
