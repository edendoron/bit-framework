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

	go reportsScheduler(10*time.Second, postIndexer)

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
