package bitHistoryCurator

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func HistoryIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-history-curator!")
}

var routes = Routes{
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
