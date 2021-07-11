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
		//TODO: handle error
		srv.Shutdown(context.Background())
	}()

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}