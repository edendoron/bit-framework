package main

import (
	storage "github.com/edendoron/bit-framework/internal/bitStorageAccess"
	"github.com/edendoron/bit-framework/internal/models"
	"github.com/edendoron/bit-framework/server"
	"log"
)

func main() {

	storage.LoadConfigs()

	models.RedirectLogger(storage.Configs.BitStoragePath)

	log.Printf("Server started - bit-storage-access")

	router := storage.StorageAccessRoutes.NewRouter()

	srv := server.NewServer(router, storage.Configs.Host+storage.Configs.BitStoragePort)

	if storage.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(storage.Configs.SSHCertPath, storage.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
