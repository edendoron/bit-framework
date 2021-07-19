package bitHandler

import (
	"fmt"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"strings"
)

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello bit-handler!")
	if err != nil {
		log.Printf("error in index route: %v", err)
	}
}

var HandlerRoutes = models.Routes{
	models.Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: HandlerIndex,
	},

	models.Route{
		Name:        "HandlerGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: HandlerGetPing,
	},

	models.Route{
		Name:        "GetTrigger",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/logic",
		HandlerFunc: GetTrigger,
	},

	models.Route{
		Name:        "PostTrigger",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/logic",
		HandlerFunc: PostTrigger,
	},

	models.Route{
		Name:        "PutResetIndications",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/reset",
		HandlerFunc: PutResetIndications,
	},
}
