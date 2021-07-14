package bitHandler

import (
	. "../models"
	"fmt"
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

var HandlerRoutes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: HandlerIndex,
	},

	Route{
		Name:        "HandlerGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: HandlerGetPing,
	},

	Route{
		Name:        "GetTrigger",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/logic",
		HandlerFunc: GetTrigger,
	},

	Route{
		Name:        "PostTrigger",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/logic",
		HandlerFunc: PostTrigger,
	},

	Route{
		Name:        "PutResetIndications",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/reset",
		HandlerFunc: PutResetIndications,
	},
}
