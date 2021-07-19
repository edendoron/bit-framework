package main

import (
	query "github.com/edendoron/bit-framework/internal/bitQuery"
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/edendoron/bit-framework/server"
	"log"
)

func main() {

	query.LoadConfigs()

	RedirectLogger(query.Configs.BitQueryPath)

	log.Printf("Server started - bit-query")

	router := query.QueryRoutes.NewRouter()

	srv := server.NewServer(router, query.Configs.Host+query.Configs.BitQueryPort)

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
