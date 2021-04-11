package apiResponseHandlers

import (
	. "../models"
	"encoding/json"
	"log"
	"net/http"
)

func ResponseHandler(w http.ResponseWriter, message string) {
	response := ApiResponse{Code: 200, Message: message}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
