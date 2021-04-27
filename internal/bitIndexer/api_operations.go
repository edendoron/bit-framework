/*
 * bitIndexer
 *
 * This protocol defines the API for **indexer** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package bitIndexer

import (
	"net/http"
)

const storageDataWriteURL = "http://localhost:8082/data/write"

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

	http.Post(storageDataWriteURL, "application/json; charset=UTF-8", r.Body)
	// ApiResponseHandler(w, resp.StatusCode, "Report posted!", nil)

}
