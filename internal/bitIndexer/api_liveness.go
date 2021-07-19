package bitindexer

import (
	"encoding/json"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"time"
)

func IndexerGetPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := models.PongBody{Timestamp: time.Now(), Version: Configs.Version, Host: Configs.Host, Ready: true, ApiVersion: Configs.ApiVersion}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalln(err)
	}
	w.WriteHeader(http.StatusOK)
}
