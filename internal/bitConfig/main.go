package bitConfig

import (
	"../../server"
	"context"
	"log"
)

func BitConfig() {
	log.Printf("Service started - bit-config")

	router := routes.NewRouter()

	srv := server.NewServer(router, ":8084")

	go func() {
		PostFailuresData()
		log.Printf("Service ended - bit-config")
		//TODO: handle error
		srv.Shutdown(context.Background())
	}()

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
