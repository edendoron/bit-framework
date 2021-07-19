package bitHandler

import (
	"encoding/json"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"time"
)

func HandlerGetPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := models.PongBody{Timestamp: time.Now(), Version: Configs.Version, Host: Configs.Host, Ready: true, ApiVersion: Configs.ApiVersion}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
