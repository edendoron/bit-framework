package bitQuery

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func QueryIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-query service!")
}

var QueryRoutes = Routes{
	Route{
		Name:        "QueryIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: QueryIndex,
	},

	Route{
		Name:        "QueryGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: QueryGetPing,
	},

	Route{
		Name:        "BitStatusQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/status",
		HandlerFunc: BitStatusQuery,
	},

	Route{
		Name:        "ReportQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/reports",
		HandlerFunc: ReportQuery,
	},
}
