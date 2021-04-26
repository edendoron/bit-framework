package bitExporter

import (
	"../../server"
	"container/heap"
	"log"
	"time"
)

func BitExporter() {

	log.Printf("Server started - bit-test-result-exporter")

	router := routes.NewRouter()

	heap.Init(&reportsQueue)

	srv := server.NewServer(router, ":8080")

	//go postIndexer()
	//
	go updateReportToIndexer()
	//
	go reportsScheduler(5 * time.Second)

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("localhost.crt", "localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
