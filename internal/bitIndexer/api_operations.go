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
	//client := &http.Client{}
	//req, err := http.NewRequest(http.MethodPut, storageDataWriteURL, r.Body)
	//if err != nil{
	//	ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	//}
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//resp, err := client.Do(req)
	//if err != nil{
	//	ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	//}
	//if resp.StatusCode == http.StatusInternalServerError{
	//	ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	//}
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)
	message := KeyValue{Key: "reports", Value: buf.String()}
	bodyRef, _ := json.MarshalIndent(message, "", " ")
	body := bytes.NewBuffer(bodyRef)
	resp, err := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", body)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
	ApiResponseHandler(w, resp.StatusCode, "Report sent to storage!", nil)
}
