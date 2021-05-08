/*
 * bit-storage-access
 *
 * This protocol defines the API for **storage-access** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package bitStorageAccess

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func StorageIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-storage-access!")
}

var StorageAccessRoutes = Routes{
	Route{
		Name:        "StorageIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: StorageIndex,
	},

	Route{
		Name:        "GetExtendedStatus",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/extended/status",
		HandlerFunc: GetExtendedStatus,
	},

	Route{
		Name:        "StorageGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: StorageGetPing,
	},

	Route{
		Name:        "GetDataRead",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/data/read",
		HandlerFunc: GetDataRead,
	},

	Route{
		Name:        "PostDataWrite",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/data/write",
		HandlerFunc: PostDataWrite,
	},

	Route{
		Name:        "PutDataRead",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/data/read",
		HandlerFunc: PutDataRead,
	},

	Route{
		Name:        "PutDataWrite",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/data/write",
		HandlerFunc: PutDataWrite,
	},
}