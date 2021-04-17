package bitIndexer

import (
	"../../server"
	. "../models"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func BitIndexer() {

	log.Printf("Server started - bit-indexer")

	router := routes.NewRouter()

	srv := server.NewServer(router, ":8081")

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bitIndexer!")
}

var routes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: index,
	},

	Route{
		Name:        "IndexerGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: IndexerGetPing,
	},

	Route{
		Name:        "IndexerPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: IndexerPostReport,
	},
}
