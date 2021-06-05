package bitHandler

import (
	. "../models"
	"fmt"
	"net/http"
	"strings"
)

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bit-handler!")
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
