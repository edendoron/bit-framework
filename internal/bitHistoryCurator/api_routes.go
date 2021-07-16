package bitHistoryCurator

import (
	"fmt"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func HistoryIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-history-curator!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var HistoryCuratorRoutes = Routes{
	Route{
		Name:        "HistoryIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: HistoryIndex,
	},
	Route{
		Name:        "StorageGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: HistoryCuratorGetPing,
	},
}
