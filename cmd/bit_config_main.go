package main

import (
	config "../internal/bitConfig"
	"../server"
	"context"
	"log"
)

func main() {
	log.Printf("Service started - bit-config")

	router := config.ConfigRoutes.NewRouter()

	srv := server.NewServer(router, ":8084")

	go func() {
		config.PostFailuresData()
		config.PostGroupFilterData()
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
