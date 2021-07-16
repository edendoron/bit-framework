package bitStorageAccess

import (
	"fmt"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func StorageIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-storage-access!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
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

	Route{
		Name:        "DeleteData",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/data/delete",
		HandlerFunc: DeleteData,
	},
}
