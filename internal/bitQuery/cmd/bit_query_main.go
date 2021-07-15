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

	if query.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(query.Configs.SSHCertPath, query.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
