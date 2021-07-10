package main

import (
	indexer ".."
	"../../../server"
	. "../../models"
	"log"
)

func main() {

	indexer.LoadConfigs()

	RedirectLogger(indexer.Configs.BitIndexerPath)

	log.Printf("Server started - bit-indexer")

	router := indexer.IndexerRoutes.NewRouter()

	srv := server.NewServer(router, indexer.Configs.BitIndexerPort)

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("../localhost.crt", "../localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
