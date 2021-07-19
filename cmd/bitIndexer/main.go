package main

import (
	indexer "github.com/edendoron/bit-framework/internal/bitIndexer"
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/edendoron/bit-framework/server"
	"log"
)

func main() {

	indexer.LoadConfigs()

	models.RedirectLogger(indexer.Configs.BitIndexerPath)

	log.Printf("Server started - bit-indexer")

	router := indexer.IndexerRoutes.NewRouter()

	srv := server.NewServer(router, indexer.Configs.Host+indexer.Configs.BitIndexerPort)

	if indexer.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(indexer.Configs.SSHCertPath, indexer.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
