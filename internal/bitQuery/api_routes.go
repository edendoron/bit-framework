package bitquery

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
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

var QueryRoutes = models.Routes{
	models.Route{
		Name:        "QueryIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: QueryIndex,
	},

	models.Route{
		Name:        "QueryGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: QueryGetPing,
	},

	models.Route{
		Name:        "BitStatusQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/status",
		HandlerFunc: BitStatusQuery,
	},

	models.Route{
		Name:        "ReportQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/reports",
		HandlerFunc: ReportQuery,
	},

	models.Route{
		Name:        "UserGroupQuery",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/userGroups",
		HandlerFunc: UserGroupQuery,
	},
}
