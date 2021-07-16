package bitQuery

import (
	"fmt"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func QueryIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-query service!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
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

	Route{
		Name:        "UserGroupQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/userGroups",
		HandlerFunc: UserGroupQuery,
	},
}
