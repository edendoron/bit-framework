package bitstorageaccess

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
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

var StorageAccessRoutes = models.Routes{
	models.Route{
		Name:        "StorageIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: StorageIndex,
	},

	models.Route{
		Name:        "GetExtendedStatus",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/extended/status",
		HandlerFunc: GetExtendedStatus,
	},

	models.Route{
		Name:        "StorageGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: StorageGetPing,
	},

	models.Route{
		Name:        "GetDataRead",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/data/read",
		HandlerFunc: GetDataRead,
	},

	models.Route{
		Name:        "PostDataWrite",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/data/write",
		HandlerFunc: PostDataWrite,
	},

	models.Route{
		Name:        "PutDataRead",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/data/read",
		HandlerFunc: PutDataRead,
	},

	models.Route{
		Name:        "PutDataWrite",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/data/write",
		HandlerFunc: PutDataWrite,
	},

	models.Route{
		Name:        "DeleteData",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/data/delete",
		HandlerFunc: DeleteData,
	},
}
