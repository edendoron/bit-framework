package bitConfig

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func ConfigIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-history-curator!")
}

var routes = Routes{
	Route{
		Name:        "ConfigIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ConfigIndex,
	},
	Route{
		Name:        "StorageGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ConfigGetPing,
	},
}
