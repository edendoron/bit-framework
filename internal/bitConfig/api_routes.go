package bitconfig

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func ConfigIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-config!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var ConfigRoutes = models.Routes{
	models.Route{
		Name:        "ConfigIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ConfigIndex,
	},
	models.Route{
		Name:        "ConfigGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ConfigGetPing,
	},
}
