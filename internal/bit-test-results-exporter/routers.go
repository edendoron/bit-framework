/*
 * bit-test-results-exporter
 *
 * This protocol defines the API for **test-results-exporter** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package bitTestResultsExporter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bitTestResultsExporter!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GetPing",
		strings.ToUpper("Get"),
		"/ping",
		GetPing,
	},

	Route{
		"GetBandwidth",
		strings.ToUpper("Get"),
		"/qos/bandwidth",
		GetBandwidth,
	},

	Route{
		"PostBandwidth",
		strings.ToUpper("Post"),
		"/qos/bandwidth",
		PostBandwidth,
	},

	Route{
		"PostReport",
		strings.ToUpper("Post"),
		"/report/raw",
		PostReport,
	},
}
