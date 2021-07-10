package bitExporter

import (
	. "../models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func ExporterGetPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := PongBody{Timestamp: time.Now(), Version: "1.00", Host: "EdenNadav", Ready: true, ApiVersion: "alpha"}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalln(err)
	}
	w.WriteHeader(http.StatusOK)
}
