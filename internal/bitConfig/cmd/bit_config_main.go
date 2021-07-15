package main

import (
	config ".."
	"../../../server"
	. "../../models"
	"context"
	"log"
)

func main() {

	config.LoadConfigs()

	RedirectLogger(config.Configs.BitConfigPath)

	log.Printf("Service started - bit-config " + config.Configs.BitConfigPort)

	router := config.ConfigRoutes.NewRouter()

	srv := server.NewServer(router, config.Configs.BitConfigPort)

	go func() {
		config.PostFailuresData()
		config.PostGroupFilterData()
		log.Printf("Service ended - bit-config")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()

	if config.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(config.Configs.SSHCertPath, config.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
