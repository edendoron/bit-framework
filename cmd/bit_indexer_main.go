package main

import (
	indexer "../internal/bitIndexer"
	. "../internal/models"
	"../server"
	"log"
)

func main() {

	RedirectLogger("./internal/bitIndexer")

	log.Printf("Server started - bit-indexer")

	router := indexer.IndexerRoutes.NewRouter()

	srv := server.NewServer(router, ":8081")

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("../localhost.crt", "../localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
