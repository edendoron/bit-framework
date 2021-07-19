package bitindexer

import (
	"bytes"
	"encoding/json"
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	"github.com/edendoron/bit-framework/internal/models"
	"io"
	"net/http"
	"strings"
)

func IndexerPostReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf := new(strings.Builder)
	n, err := io.Copy(buf, r.Body)
	if err != nil || n != r.ContentLength {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	message := models.KeyValue{Key: "reports", Value: buf.String()}
	bodyRef, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	body := bytes.NewBuffer(bodyRef)
	resp, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", body)
	if err != nil || resp.StatusCode != http.StatusOK {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer resp.Body.Close()
	rh.ApiResponseHandler(w, http.StatusOK, "Report sent to storage!", nil)
}
