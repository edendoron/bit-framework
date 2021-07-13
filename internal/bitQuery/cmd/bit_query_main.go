package main

import (
	query ".."
	"../../../server"
	. "../../models"
	"log"
)

func main() {

	query.LoadConfigs()

	RedirectLogger(query.Configs.BitQueryPath)

	log.Printf("Server started - bit-query")

	router := query.QueryRoutes.NewRouter()

	srv := server.NewServer(router, query.Configs.BitQueryPort)

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS(query.Configs.SSHCertPath, query.Configs.SSHKeyPath)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
