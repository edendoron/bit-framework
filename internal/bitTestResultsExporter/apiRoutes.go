package bitExporter

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func ExporterIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bitTestResultsExporter!")
}

var routes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ExporterIndex,
	},

	Route{
		Name:        "ExporterGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ExporterGetPing,
	},

	Route{
		Name:        "GetBandwidth",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: GetBandwidth,
	},

	Route{
		Name:        "PostBandwidth",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: PostBandwidth,
	},

	Route{
		Name:        "ExporterPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: ExporterPostReport,
	},
}
