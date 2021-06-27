package main

import (
	query "../internal/bitQuery"
	. "../internal/models"
	"../server"
	"log"
)

func main() {

	RedirectLogger("./internal/bitQuery")

	log.Printf("Server started - bit-query")

	router := query.QueryRoutes.NewRouter()

	srv := server.NewServer(router, ":8085")

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("../localhost.crt", "../localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
