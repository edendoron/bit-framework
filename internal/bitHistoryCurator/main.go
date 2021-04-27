package bitHistoryCurator

import (
	"../../server"
	"log"
)

func BitHistoryCurator() {
	log.Printf("Server started - bit-history-curator")

	router := routes.NewRouter()

	srv := server.NewServer(router, ":8083")

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
