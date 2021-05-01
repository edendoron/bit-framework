package bitIndexer

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func IndexerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bitIndexer!")
}

var routes = Routes{
	Route{
		Name:        "Index",
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
