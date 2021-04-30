/*
 * bit-storage-access
 *
 * This protocol defines the API for **storage-access** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package bitStorageAccess

import (
	"../../server"
	"log"
)

func BitStorageAccess() {
	log.Printf("Server started - bit-storage-access")

	router := routes.NewRouter()

	srv := server.NewServer(router, ":8082")

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}