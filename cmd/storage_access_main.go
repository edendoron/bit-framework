package main

import (
	storage "../internal/bitStorageAccess"
	. "../internal/models"
	"../server"
	"log"
)

func main() {

	RedirectLogger("./internal/bitStorageAccess")

	log.Printf("Server started - bit-storage-access")

	router := storage.StorageAccessRoutes.NewRouter()

	srv := server.NewServer(router, ":8082")

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
