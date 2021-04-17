package bitExporter

import (
	"../../server"
	. "../models"
	"container/heap"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func BitExporter() {

	log.Printf("Server started - bit-test-result-exporter")

	router := routes.NewRouter()

	heap.Init(&reportsQueue)

	srv := server.NewServer(router, ":8080")

	go reportsScheduler(10*time.Second, postIndexer)

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bitTestResultsExporter!")
}

var routes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: index,
	},

	Route{
		Name:        "ExporterGetPing",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/ping",
		HandlerFunc: ExporterGetPing,
	},

	Route{
		Name:        "GetBandwidth",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: GetBandwidth,
	},

	Route{
		Name:        "PostBandwidth",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/qos/bandwidth",
		HandlerFunc: PostBandwidth,
	},

	Route{
		Name:        "ExporterPostReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/report/raw",
		HandlerFunc: ExporterPostReport,
	},
}
