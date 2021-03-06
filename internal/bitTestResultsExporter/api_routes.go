package bitexporter

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func ExporterIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit Test Results Exporter!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var ExporterRoutes = models.Routes{
	models.Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ExporterIndex,
	},

	models.Route{
		Name:        "ExporterGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ExporterGetPing,
	},

	models.Route{
		Name:        "GetBandwidth",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: GetBandwidth,
	},

	models.Route{
		Name:        "PostBandwidth",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: PostBandwidth,
	},

	models.Route{
		Name:        "ExporterPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: ExporterPostReport,
	},
}
