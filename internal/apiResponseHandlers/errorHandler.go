package apiResponseHandlers

import (
	. "../models"
	"encoding/json"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, message string, err error) {
	badResponse := ApiResponse{Code: 404, Message: message}
	w.WriteHeader(http.StatusBadRequest)
	e := json.NewEncoder(w).Encode(&badResponse)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(e)
	}
	log.Fatalln(err)
}
