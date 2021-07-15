package main

import (
	storage ".."
	"../../../server"
	. "../../models"
	"log"
)

func main() {

	storage.LoadConfigs()

	RedirectLogger(storage.Configs.BitStoragePath)

	log.Printf("Server started - bit-storage-access")

	router := storage.StorageAccessRoutes.NewRouter()

	srv := server.NewServer(router, storage.Configs.BitStoragePort)

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
