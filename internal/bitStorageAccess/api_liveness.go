package bitstorageaccess

import (
	"encoding/json"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"time"
)

func GetExtendedStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := models.StorageExtendedStatus{StorageType: "file_system"}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalln(err)
	}
	w.WriteHeader(http.StatusOK)
}

func StorageGetPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := models.PongBody{Timestamp: time.Now(), Version: Configs.Version, Host: Configs.Host, Ready: true, ApiVersion: Configs.ApiVersion}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalln(err)
	}
	w.WriteHeader(http.StatusOK)
}
