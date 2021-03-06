package bitindexer

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func IndexerIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-indexer!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var IndexerRoutes = models.Routes{
	models.Route{
		Name:        "IndexerIndex",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: IndexerIndex,
	},

	models.Route{
		Name:        "IndexerGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: IndexerGetPing,
	},

	models.Route{
		Name:        "IndexerPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: IndexerPostReport,
	},
}
