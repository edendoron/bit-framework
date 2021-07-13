package bitConfig

import (
	. "../models"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ConfigIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Println(w, "Hello bit-config!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var ConfigRoutes = Routes{
	Route{
		Name:        "ConfigIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ConfigIndex,
	},
	Route{
		Name:        "ConfigGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ConfigGetPing,
	},
}
