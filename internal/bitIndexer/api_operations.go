package bitIndexer

import (
	. "../apiResponseHandlers"
	. "../models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func IndexerPostReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf := new(strings.Builder)
	n, err := io.Copy(buf, r.Body)
	if err != nil || n != r.ContentLength {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	message := KeyValue{Key: "reports", Value: buf.String()}
	bodyRef, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	body := bytes.NewBuffer(bodyRef)
	resp, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", body)
	if err != nil || resp.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer resp.Body.Close()
	ApiResponseHandler(w, http.StatusOK, "Report sent to storage!", nil)
}
