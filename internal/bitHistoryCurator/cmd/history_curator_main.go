package main

import (
	history ".."
	"../../../server"
	. "../../models"
	"context"
	"log"
)

func main() {

	history.LoadConfigs()

	RedirectLogger(history.Configs.BitHistoryCuratorPath)

	log.Printf("Service started - bit-history-curator")

	router := history.HistoryCuratorRoutes.NewRouter()

	srv := server.NewServer(router, history.Configs.BitHistoryCuratorPort)

	go func() {
		history.RemoveAgedData(history.GetCuratorTimeConfig())
		log.Printf("Service ended - bit-history-curator")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	if history.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(history.Configs.SSHCertPath, history.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
