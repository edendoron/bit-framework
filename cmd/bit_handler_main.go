package main

import (
	handler "../internal/bitHandler"
	"../server"
	"log"
)

func main() {
	log.Printf("Server started - bit-storage-access")

	router := handler.HandlerRoutes.NewRouter()

	srv := server.NewServer(router, ":8085")

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
