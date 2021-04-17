package bitIndexer

import (
	"../../server"
	"log"
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
