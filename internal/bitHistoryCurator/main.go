package bitHistoryCurator

import (
	"../../server"
	"context"
	"log"
)

func BitHistoryCurator() {
	log.Printf("Service started - bit-history-curator")

	router := routes.NewRouter()

	srv := server.NewServer(router, ":8083")

	go func() {
		removeAgedData(getConfig().AgedTimeDuration)
		log.Printf("Service ended - bit-history-curator")
		//TODO: handle error
		srv.Shutdown(context.Background())
	}()

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
