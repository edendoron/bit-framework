package bitIndexer

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func IndexerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-indexer!")
}

var routes = Routes{
	Route{
		Name:        "IndexerIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: IndexerIndex,
	},

	Route{
		Name:        "IndexerGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: IndexerGetPing,
	},

	Route{
		Name:        "IndexerPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: IndexerPostReport,
	},
}
